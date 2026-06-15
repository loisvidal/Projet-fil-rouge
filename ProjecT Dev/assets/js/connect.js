const modalContainer = document.querySelector(".modal-container");
const modalTrigger = document.querySelectorAll(".modal-trigger");

if (modalTrigger.length > 0) {
    modalTrigger.forEach(trigger => trigger.addEventListener("click", toggleModal));
}

function toggleModal() {
    if (modalContainer) {
        modalContainer.classList.toggle("active");
    }
}

const login = document.querySelector(".login");
const register = document.querySelector(".register");
const switchLogin = document.querySelector(".switch-login");
const switchRegister = document.querySelector(".switch-register");

if (switchLogin) {
    switchLogin.addEventListener("click", toggleLogin);
}
if (switchRegister) {
    switchRegister.addEventListener("click", toggleRegister);
}

function toggleLogin() {
    if (register) register.classList.remove("active");
    if (login) login.classList.add("active");
}

function toggleRegister() {
    if (login) login.classList.remove("active");
    if (register) register.classList.add("active");
}
