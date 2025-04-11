
const calculator = require("./calculator");

const args = process.argv.slice(2);

const operator = args[0];

const numbers = args.slice(1).map(Number);

if(!operator || numbers.length < 2){
    console.error("请按规定输入：node index.js <操作符> <数字1> <数字2>...");
    process.exit(1);
}

try {
    let result;
    switch (operator) {
        case "+":
            result = calculator.add(...numbers);
            break;
        case "-":
            result = calculator.subtract(...numbers);
            break;
        case "*":
            result = calculator.multiply(...numbers);
            break;
        case "/":
            result = calculator.divide(...numbers);
            break;
        default:
            console.error(`无效操作符:${operator}`);
            process.exit(1); 
        }
        console.log(`${result}`);
    } catch (error) {
        console.error(`错误:${error.message}`);
}