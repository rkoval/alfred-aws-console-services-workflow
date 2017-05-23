#!/usr/bin/env node

YAML = require('yamljs');
const consoleServices = YAML.load('console-services.yml');
const listItems = generate(consoleServices);
injectIntoXml(listItems);

function generate(consoleServices) {
  return consoleServices.map(({command, description, icon, url}) => {
    console.log(command);

    const consoleService = {
      title: command,
      arg: url,
      subtitle: description,
      imagefile: icon || `${command}.png`,
    };

    return consoleService;
  });
}

function injectIntoXml(alfredProjects) {
  const plist = require('plist');
  const fs = require('fs-extra');
  const path = require('path');

  fs.copySync(path.resolve(`${__dirname}/../info.plist`), `${__dirname}/../info.plist.generatebackup`);
  const file = fs.readFileSync(`${__dirname}/../info.plist.generatebackup`, 'utf-8');
  const plistJson = plist.parse(file);
  plistJson.objects[1].config.items = JSON.stringify(alfredProjects);
  fixedPlistJson = fixPlistLibraryNullKey(plistJson);
  const updatedPlist = plist.build(fixedPlistJson)
  fs.writeFileSync(`${__dirname}/../info.plist`, updatedPlist);
}

function fixPlistLibraryNullKey(plistJson) {
  const traverse = require('traverse');
  traverse(plistJson).forEach(function(node) {
    if (node === null) {
      this.update('');
    }
  })
  return plistJson;
}
