#!/bin/bash

gometalinter | grep -v "_test\.go.*[(]lll[)]$"
