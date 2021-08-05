#!/usr/bin/env node
const fs = require('fs');
const assert = require('assert-diff');

function main() {
  const stdinBuffer = fs.readFileSync(0);
  const rawSession = stdinBuffer.toString();
  const session = JSON.parse(rawSession);
  const activeAwsUrls = [];
  session.windows.forEach(window => {
    window.tabs.forEach(tab => {
      const currentEntry = tab.entries[tab.entries.length - 1];
      if (currentEntry.url.includes('console.aws.amazon.com')) {
        activeAwsUrls.push(currentEntry.url);
      }
    });
  });

  const lastOpenedUrlsBuffer = fs.readFileSync(
    // TODO this will break with Alfred 4
    `${process.env.HOME}/Library/Caches/com.runningwithcrayons.Alfred-3/Workflow Data/com.ryankoval.awsconsoleservices/last-opened-urls.json`
  );
  const lastOpenedUrls = JSON.parse(lastOpenedUrlsBuffer);

  assert.deepStrictEqual(activeAwsUrls, lastOpenedUrls);
  console.log('✅ urls matches successfully!');
}

main();
