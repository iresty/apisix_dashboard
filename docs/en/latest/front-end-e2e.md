---
title: Frontend E2E
---

<!--
#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
-->

This project uses [Cypress](https://www.cypress.io/) as the front-end E2E test framework.

1. To start the front-end project locally, please refer to [develop](./develop.md) web section.

2. Open Cypress test-runner. For the use of test-runner, please refer to [test-runner](https://docs.cypress.io/guides/core-concepts/test-runner.html#) Overview.

   ```sh
   pnpm cypress:open-dev
   ```

3. Write your test examples: please refer to the test examples in the `/web/cypress` directory, or see [RWA](https://github.com/cypress-io/cypress-realworld-app) for more examples.

To make it easy for users to develop front-end E2E cases, we use the remote manager-api by default. If you want to use the local manager-api, please read the following instructions:

1. Start the local manager-api service, please refer to [develop](./develop.md) manager-api section.

2. To start the front-end project locally, please refer to [develop](./develop.md) web section. NOTE: You need to change `pnpm start` to `pnpm start:e2e` when you start.

3. Open Cypress test-runner.

   ```sh
   pnpm cypress:open
   ```

Reference links:

- Cypress API Documents: https://docs.cypress.io/api/api/table-of-contents.html

- Cypress Best Practices: https://docs.cypress.io/guides/references/best-practices.html#Organizing-Tests-Logging-In-Controlling-State
