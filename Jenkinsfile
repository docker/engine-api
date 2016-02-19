stage name: "test"

def jobFor(platform, go_version) {
  { it ->
    node(platform) {
      wrap([$class: 'TimestamperBuildWrapper']) {
        wrap([$class: 'AnsiColorBuildWrapper']) {
          sh 'rm -rf ./*'
          if (platform == "windows") { tool 'hg' }
          checkout scm
          withEnv([
            "GOVERSION=${go_version}",
            "GOCYCLO_MAX=21",
            "GOPACKAGE=github.com/docker/engine-api"
          ]) {
            withCredentials([[$class: 'StringBinding', credentialsId: 'docker-jenkins.token.github.com', variable: 'GITHUB_TOKEN']]) {
              try {
                sh "bash hack/test-${platform}.sh || exit 0"
              } finally {
                if (platform == "linux") {
                  sh "docker run --rm -v \"\$(pwd):/workspace\" busybox chown -R \"\$(id -u):\$(id -g)\" /workspace"
                }
              }
            }
            step([$class: 'JUnitResultArchiver', testResults: 'results/tests.xml', keepLongStdio: true])
            // step([$class: 'hudson.plugins.cobertura.CoberturaPublisher', coberturaReportFile: 'results/coverage.xml'])
            step([
              $class: 'WarningsPublisher',
              parserConfigurations: [[
                parserName: "Go Lint",
                pattern: "results/fmt.txt,results/lint.txt,results/cyclo.txt",
              ], [
                parserName: "Go Vet",
                pattern: "results/vet.txt"
              ]],
              unstableTotalAll: '0'
            ])
            archive 'results'
          }
        }
      }
    }
  }
}

parallel(
  failFast: false,
  "windows-1.5.3": jobFor("windows", "1.5.3"),
  "linux-1.5.3": jobFor("linux", "1.5.3"),
  "linux-1.6": jobFor("linux", "1.6")
)
