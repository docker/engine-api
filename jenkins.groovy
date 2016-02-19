import groovy.json.JsonSlurper

// Change this as complexity of functions is reduced
// The value is the max allowed which will still pass the job.
def GOCYCLO_MAX = 21

// We want to avoid hard-coding the owner and repo name so that this
// file stands on its own if this project is forked and someone wants to
// use this DSL script.

// Step 1: find out where this project came from.
//   This is really ugly but since we don't have access to the workspace
//   it's the best I could figure out. Things to be aware of:
//   * jm.delegate is because the jm class is javaposse.jobdsl.plugin.InterruptibleJobManagement.
//      if it ever changes back to JenkinsJobManagement (or other) we'll need to update that part.
//   * if scm is not configured or not git, there's going to be a problem.
//   * if scm has more than one remote configured, the first one is used.
def scmUrl = jm.delegate.build.project.scm.userRemoteConfigs[0].url

// Step 2: Extract the important data github URL
//   This logic should work with these variants:
//   * [https|git]://[www.]github.com/owner/name[.git]
//   * [ssh://]git@github.com:owner/name[.git]
def (owner, repoName) = scmUrl.split(":")[-1].split("/")[-2..-1]
if (repoName.endsWith(".git")) {
  repoName = repoName[0..-5]
}
def projectPrefix = "${owner}-${repoName}"

// Step 3: collect some information about the repo from github API.
//   We want the default branch and the info about the upstream if we're a fork.
//   NOTE: this only works with public projects. org.kohsuke.github.GitHub
//   should be used for private repos and anything more advanced.
def repoInfo = new JsonSlurper().parseText(new URL("https://api.github.com/repos/${owner}/${repoName}").text)
def defaultBranch = repoInfo.default_branch
def upstreamOwner = owner
def upstreamName = repoName
if (repoInfo.fork) {
  (upstreamOwner, upstreamName) = repoInfo.parent.full_name.split("/")
}

// We'll want branches and PR heads. Tags are pulled by default
Closure refSpecParts = { remoteName ->
  "+refs/heads/*:refs/remotes/${remoteName}/* +refs/pull/*/head:refs/remotes/${remoteName}/pr/*"
}

// Instead of duplicating these options below, create a closure
// which can be applied when necessary.
Closure defaultGitSettingsClosure = { gitContext ->
  gitContext.with {
    wipeOutWorkspace()

    // If the repo were private we'd also need credentials here.
    remote {
      name("origin")
      github("${owner}/${repoName}")
      refspec(refSpecParts("origin"))
    }
  }
}

// Some maps to keep track of things
def taskJobs = [:]
def triggerJobs = [:]

folder("${projectPrefix}")

// Some jobs use volumes in the workspace and end up with the files owned by root.
// This function puts a chown trap on the first line of a script to fix ownership.
def chownWorkspaceShell(cmd) {
  def lines = cmd.readLines()
  def doChown = "trap 'docker run --rm -v \"\$(pwd):/workspace\" busybox chown -R \"\$(id -u):\$(id -g)\" /workspace' EXIT"
  def index = 0
  if (lines[0].startsWith("#!")) {
    index = 1
  }
  lines.add(index, doChown)
  return lines.join("\n")
}

// TODO: use matrixJob for windows as well
taskJobs["windows"] = job("${projectPrefix}/task-test-windows") {}
taskJobs["linux"] = matrixJob("${projectPrefix}/task-test-linux") {}

// Common functionality between windows and linux.
taskJobs.values().each {
  it.with {
    wrappers {
      environmentVariables {
        env('GOCYCLO_MAX', GOCYCLO_MAX)
        env('GOPACKAGE', "github.com/${upstreamOwner}/${upstreamName}")
      }
    }

    publishers {
      // Not really necessary to archive all the raw files, but just in case someone
      // wants to look at them
      archiveArtifacts('results/**')
      archiveJunit("results/tests.xml") { retainLongStdout() }
      cobertura("results/coverage.xml")
      warnings([], [
          "Go Lint": "results/fmt.txt,results/lint.txt,results/cyclo.txt",
          "Go Vet": "results/vet.txt"
      ]) {
        thresholds(unstableTotal: [all: 0])
      }
    }

    // NOTE: these steps also be moved to "main" if users are OK with browsing around jenkins a bit.
    // It basically makes it less cluttered for success cases, but harder to find problems.
    // Since problems are more important and I don't think anyone minds a string of green successes,
    // I opted for this way.
    // TODO: implement commit status notifier in DSL? node buildier is awkward and error-prone.
    configure { node ->
      node / builders / "com.cloudbees.jenkins.GitHubSetCommitStatusBuilder" / "statusMessage" {
        content("")
      }
      node / "publishers" / "com.cloudbees.jenkins.GitHubCommitNotifier" {
        delegate / statusMessage / "content" ("")
        delegate / resultOnFailure("FAILURE")
      }
    }
  }
}

taskJobs["linux"].with {
  axes {
    labelExpression("label", ["docker && linux"])
    text("GOVERSION", ["1.5.3", "1.6"])
  }

  steps {
    shell(chownWorkspaceShell(readFileFromWorkspace("hack/test-linux.sh")))
  }
}

// TODO: use matrixJob and the axes below
taskJobs["windows"].with {
  // axes {
  //   labelExpression("label", ["windows"])
  //   // text("GOVERSION", ["1.5.3"])
  // }

  label("windows")

  steps {
    shell(readFileFromWorkspace("hack/test-windows.sh"))
  }

  wrappers {
    credentialsBinding {
      // mercurial is necessary to install some of the go tools
      customTools(["hg"]) { skipMasterInstallation() }
      // Token is needed to pull install script from private github repo
      string("GITHUB_TOKEN", "docker-jenkins.token.github.com")
    }
  }
}

def main
main = multiJob("${projectPrefix}/main") {
  steps {
    phase("test", "UNSTABLE") {
      taskJobs.values().each {
        phaseJob(it.name) {
          parameters {
            predefinedProps(GIT_REF: '$GIT_COMMIT')
          }
        }
      }
    }
  }
}

// tests and 'main' all accept a git ref
(taskJobs.values() + [main]).each {
  it.with {
    parameters {
      stringParam("GIT_REF", "origin/${defaultBranch}", "What to build. Anything git understands is acceptable.")
    }
    scm {
      git {
        // Non-public projects or jobs that push will need credentials.
        defaultGitSettingsClosure(delegate)
        branch('$GIT_REF')
      }
    }
  }
}

// All the different ways to trigger from our repo.
// Pulls need different trigger for GHPRB plugin
triggerJobs[defaultBranch] = multiJob("${projectPrefix}/${defaultBranch}") {}
triggerJobs['branches'] = multiJob("${projectPrefix}/branches") {}
triggerJobs['tags'] = multiJob("${projectPrefix}/tags") {}
triggerJobs['pulls'] = multiJob("${projectPrefix}/pulls") {
  triggers {
    // TODO: pullRequest says it's deprecated but what's the replacement?
    pullRequest {
      useGitHubHooks()
      orgWhitelist([upstreamOwner, owner])
      allowMembersOfWhitelistedOrgsAsAdmin()
      permitAll()
    }
  }
}

triggerJobs.each { k, triggerJob ->
  if (k == "branches" || k == defaultBranch) {
    refToBuild = "origin/${defaultBranch}"
  } else if (k == "pull") {
    refToBuild = '${ghprbActualCommit}'
  } else if (k == "tags") {
    refToBuild = "tags/*"
  }

  triggerJob.with {
    triggers { githubPush() }
    scm {
      git {
        defaultGitSettingsClosure(delegate)
        branch(refToBuild)
        if (k == "branches") {
          // This will build everything *except* the default branch
          strategy { inverse() }
        }
      }
    }
    steps {
      phase("main", "UNSTABLE") {
        phaseJob(main.name) {
          parameters {
            predefinedProps(GIT_REF: '$GIT_COMMIT')
          }
        }
      }
    }
  }
}

// Settings for all the jobs we've created
(taskJobs.values() + triggerJobs.values() + [main]).each {
  it.with {
    logRotator { numToKeep(250) }
    concurrentBuild()
    wrappers {
      timestamps()
      colorizeOutput()
    }
    properties {
      ownership {
        // TODO: How to automatically detect owners?
        primaryOwnerId("miked")
      }
    }
  }
}
