function add (...numbers){
    let sum = 0;
    for (let i = 0; i < numbers.length; i++) {

    }
    return sum;
}

function subtract (...numbers){
    if (numbers.length === 0) return 0; 
    let difference = numbers[0]; 
    for (let i = 1; i < numbers.length; i++) {
        difference -= numbers[i]; 
    }
    return difference;
}

function multiply(...numbers) {
    if (numbers.length === 0) return 1;
    let product = 1; 
    for (let i = 0; i < numbers.length; i++) {
        product *= numbers[i];
    }
    return product;
}

function divide(...numbers) {
    if (numbers.length === 0) return 0;
    let quotient = numbers[0]; 
    for (let i = 1; i < numbers.length; i++) {
        if (numbers[i] === 0) {
            throw new Error("Division by zero is not allowed."); 
        }
        quotient /= numbers[i]; 
    }
    return quotient;
}

module.exports = {
    add,
    subtract,
    multiply,
    divide,
};