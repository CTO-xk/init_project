const chalk = require("chalk");
/**
 * 日志工具
 * 支持不同级别日志和模块区分
 */
function getLogger(moduleName) {
  const prefix = `[${new Date().toISOString()}] [${moduleName}]`;
  return {
    info: (message) => {
      console.log(chalk.blue(`${prefix} INFO: ${message}`));
    },
    success: (message) => {
      console.log(chalk.green(`${prefix} SUCCESS: ${message}`));
    },
    warn: (message) => {
      console.log(chalk.yellow(`${prefix} WARN: ${message}`));
    },
    error: (message) => {
      console.error(chalk.red(`${prefix} ERROR: ${message}`));
    },
  };
}

module.exports = { getLogger };