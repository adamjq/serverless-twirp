#!/usr/bin/env node
import 'source-map-support/register';
import { App } from '@aws-cdk/core';
import { ApiStack } from '../lib/api';

const app = new App();
new ApiStack(app, 'ServerlessTwirpApiStack', {});
