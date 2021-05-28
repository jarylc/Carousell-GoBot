#!/usr/bin/env sh
npm ci
npx browserify chrono.js --standalone chrono > chrono.out.js
go-bindata -pkg chrono chrono.out.js
rm -f chrono.out.js
