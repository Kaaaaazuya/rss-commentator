#!/usr/bin/env node
import * as cdk from 'aws-cdk-lib';
import { RssServiceStack } from '../lib/cdk-stack';

const app = new cdk.App();
new RssServiceStack(app, 'RssService', {});
