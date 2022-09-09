# Testing

The definitions and distinction between testing levels could be a bit vague.
As

### Unit tests 

Maximum coverage, part of build pipeline. Test targets are functions.

* Package level
* Use _test suffix for test package to avoid testing implementation details
* Use interfaces to inject dependencies
* Use https://github.com/golang/mock to mock interfaces

### Integration tests

Separate pipeline. Cover functions. Can be part of the build pipeline or triggered by the build pipeline.
* dockertest/docker-compose
* Mountbank

### E2E/manual tests

Background checks that are not a part of the CI/CD pipeline. Serve for monitoring purposes. Do not block development iterations. Done by scheduled jobs. (Selenium, Cypress, etc…)
