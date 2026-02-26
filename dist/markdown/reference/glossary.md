
# Glossary

What we talk about when we talk about k6.

In discussion about k6, some terms have a precise, technical meaning.
If a certain term in these docs confuses you, consult this list for a definition.

- [Application performance monitoring](#application-performance-monitoring)
- [Concurrent sessions](#concurrent-sessions)
- [Checks](#checks)
- [Custom resource](#custom-resource)
- [Data correlation](#data-correlation)
- [Data parameterization](#data-parameterization)
- [Dynamic data](#dynamic-data)
- [Endurance testing](#endurance-testing)
- [Environment variables](#environment-variables)
- [Execution segment](#execution-segment)
- [Sobek](#sobek)
- [Graceful stop](#graceful-stop)
- [Happy path](#happy-path)
- [HTTP archive](#http-archive)
- [Iteration](#iteration)
- [k6 Cloud](#k6-cloud)
- [k6 options](#k6-options)
- [Kubernetes](#kubernetes)
- [Load test](#load-test)
- [Load zone](#load-zone)
- [Lifecycle function](#lifecycle-function)
- [Metric](#metric)
- [Metric sample](#metric-sample)
- [Operator pattern](#operator-pattern)
- [Parallelism](#parallelism)
- [Reliability](#reliability)
- [Requests per second](#requests-per-second)
- [Saturation](#saturation)
- [Scenario](#scenario)
- [Scenario executor](#scenario-executor)
- [Smoke test](#smoke-test)
- [Soak test](#soak-test)
- [Stability](#stability)
- [Stress test](#stress-test)
- [System under test](#system-under-test)
- [Test run](#test-run)
- [Test concurrency](#test-concurrency)
- [Test duration](#test-duration)
- [Test script](#test-script)
- [Threshold](#threshold)
- [Throughput](#throughput)
- [Virtual user](#virtual-user)
- [YAML](#yaml)
  

## Application performance monitoring

_(Or APM)_. The practice of monitoring the performance, availability, and reliability of a system.

## Concurrent sessions

The number of simultaneous VU requests in a test run.

## Checks

Checks are true/false conditions that evaluate the content of some value in the JavaScript runtime.Checks reference

## Custom resource

An extension to the [Kubernetes](#kubernetes) API.[Kubernetes reference](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)

## Data correlation

The process of taking [dynamic data](#dynamic-data) received from the system under test and reusing the data in a subsequent request.

Correlation and dynamic data example,Correlation in testing APIs

## Data parameterization

The process of turning test values into reusable parameters, e.g. through variables and shared arrays.Data parameterization examples

## Dynamic data

Data that might change or will change during test runs or across test runs. Common examples are order IDs, session tokens, or timestamps.Correlation and dynamic data example

## Endurance testing

A synonym for [soak testing](#soak-test).

## Environment variables

User-definable values which may be utilized by the operating system and other programs.Using environment variables

## Execution segment

A partition, or fractional portion, of an overall [test run](#test-run).The execution-segment options

## Sobek

A JavaScript engine written in Go. k6 binaries are embedded with Sobek, enabling test scripting in JavaScript.[Sobek repository](https://github.com/grafana/sobek) fork of [goja](https://github.com/dop251/goja).

## Graceful stop

A period that lets VUs finish an iteration at the end of a load test. Graceful stops prevent abrupt halts in execution.Graceful stop reference

## Happy path

The default system behavior that happens when a known input produces an expected output. A common mistake in performance testing happens when scripts account only for the best case (in other words, the happy path). Most load tests try to discover system errors, so test scripts should include exception handling.[Happy path (Wikipedia)](https://en.wikipedia.org/wiki/Happy_path)

## HTTP archive

_(Or HAR file)_. A file containing logs of browser interactions with the system under test. All included transactions are stored as JSON-formatted text. You can use these archives to generate test scripts (for example, with the har-to-k6 Converter).[HAR 1.2 Specification](http://www.softwareishard.com/blog/har-12-spec/),HAR converter

## Iteration

A single run in the execution of the `default function`, or scenario `exec` function. You can set iterations across all VUs, or per VU.The test life cycle document breaks down each stage of a k6 script, including iterations in VU code.

## JSON

An open standard, human-readable data-serialization format originally derived from JavaScript.

## k6 Cloud

The proper name for the entire cloud product, comprising both k6 Cloud Execution and k6 Cloud Test Results.[k6 Cloud docs](https://grafana.com/docs/grafana-cloud/testing/k6/)

## k6 options

Values that configure a k6 test run. You can set options with command-line flags, environment variables, and in the script.k6 Options

## Kubernetes

An open-source system for automating the deployment, scaling, and management of containerized applications.[Kubernetes website](https://kubernetes.io/)

## Load test

A test that assesses the performance of the system under test in terms of concurrent users or requests per second.Load Testing

## Load zone

The geographical instance from which a test runs.[Private load zones](https://grafana.com/docs/grafana-cloud/testing/k6/author-run/private-load-zone-v2/),[Declare load zones from the CLI](https://grafana.com/docs/grafana-cloud/testing/k6/author-run/use-load-zones/)

## Lifecycle function

A function called in a specific sequence in the k6 runtime. The most important lifecycle function is the default function, which runs the VU code.Test lifecycle

## Metric

A measure of how the system performs during a test run. `http_req_duration` is an example of a built-in k6 metric. Besides built-ins, you can also create custom metrics.Metrics

## Metric sample

A single value for a metric in a test run. For example, the value of `http_req_duration` from a single VU request.

## Operator pattern

Extends [Kubernetes](#kubernetes), enabling cluster management of [custom resources](#custom-resource).[Kubernetes reference](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

## Parallelism

The simultaneous execution of multiple tasks by dividing a problem into smaller independent parts.

## Reliability

The probability that a system under test performs as intended.

## Requests per second

The rate at which a test sends requests to the system under test.

## Saturation

A condition when a system's reaches full resource utilization and can handle no additional request.

## Scenario

An object in a test script that makes in-depth configurations to how VUs and iterations are scheduled. With scenarios, your test runs can model diverse traffic patterns.Scenarios reference

## Scenario executor

An property of a [scenario](#scenario) that configures VU behavior.

You can use executors to configure whether to designate iterations as shared between VUs or to run per VU, or to configure or whether the VU concurrency is constant or changing.Executor reference

## Smoke test

A regular load test configured for minimum load. Smoke tests verify that the script has no errors and that the system under test can handle a minimal amount of load.Smoke Testing

## Soak test

A test that tries to uncover performance and reliability issues stemming from a system being under pressure for an extended duration.Soak Testing

## Stability

A system under test’s ability to withstand failures and errors.

## Stress test

A test that assess the availability and stability of the system under heavy load.Stress Testing

## System under test

The software that the load test tests. This could be an API, a website, infrastructure, or any combination of these.

## Test run

An individual execution of a test script over all configured iterations.[Running k6](https://grafana.com/docs/k6/v1.5.0/get-started/running-k6)

## Test concurrency

In k6 Cloud, the number of tests running at the same time.

## Test duration

The length of time that a test runs. When duration is set as an option, VU code runs for as many iterations as possible in the length of time specified.Duration option reference

## Test script

The actual code that defines how the test behaves and what requests it makes, along with all (or at least most) configuration needed to run the test.Single Request example.

## Threshold

A pass/fail criteria that evaluates whether a metric reaches a certain value. Testers often use thresholds to codify SLOs.Threshold reference

## Throughput

The rate of successful message delivery. In k6, throughput is measured in requests per second.

## Virtual user

_(Or VU)_. The simulated users that run separate and concurrent iterations of your test script.The VU option

## YAML

Rhymes with "camel," provides a human-readable data-serialization format commonly used for configuration files.

