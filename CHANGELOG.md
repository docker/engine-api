# Changelog

Items starting with DEPRECATE are important deprecation notices. For more information on the list of deprecated APIs please have a look at https://docs.docker.com/misc/deprecated/ where target removal dates can also be found.

## 0.1.1 (2016-01-06)

### Client

- Delegate shmSize units conversion to the consumer.

### Types

- Add warnings to the volume list reponse.
- Fix image build options:
	* use 0 as default value for shmSize.

## 0.1.0 (2016-01-04)

### Client

- Initial API client implementation.

### Types

- Initial API types implementation.
