/**
 * @fileoverview read our config.json files, setting default values
 */

function titleCase(s) {
  return s.charAt(0).toUpperCase() + s.substr(1);
}

function readConfig() {
  const config = require("./rulesets.json");
  config.rulesets.forEach((cfg) => {
    cfg.shortname ||= titleCase(
      cfg.ghrepo.split("/")[1].substring("rules_".length)
    );
    cfg.repository ||= cfg.ghrepo.split("/")[1].replace("-", "_");
  });
  return config;
}

module.exports = readConfig;

if (require.main === module) {
  console.log(JSON.stringify(readConfig()));
}
