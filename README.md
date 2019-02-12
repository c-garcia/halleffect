# hall effect

Lambda to poll [Concourse metrics](https://concourse-ci.org/) and publish them to AWS CloudWatch.

The main use case for it are those teams using a hosted Concourse 
instance that they cannot manage (except for creating pipelines) and
wanting to get some performance metrics how the delivery
process is going.

At the moment, it is only able to publish to AWS Cloudwatch Metrics under 
the namespace _Concourse/Jobs_.

# How to use it

The `Makefile` and the `terraform` code should provide a hint until 
there is some decent documentation available.

The lambda needs to have `CONCOURSE_NAME` and `CONCOURSE_URL` environment 
variables defined.

# License

Copyright 2019 Cristóbal García

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

