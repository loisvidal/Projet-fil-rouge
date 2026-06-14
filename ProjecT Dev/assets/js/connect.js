const modalContainer = document.querySelector(".modal-container");
const modalTrigger = document.querySelectorAll(".modal-trigger");

modalTrigger.forEach(trigger => trigger.addEventListener("click", toggleModal))

function toggleModal() {
    modalContainer.classList.toggle("active");
}

//

const login = document.querySelector(".login");
const register = document.querySelector(".register");
const switchLogin = document.querySelector(".switch-login");
const switchRegister = document.querySelector(".switch-register");

switchLogin.addEventListener("click", toggleLogin);
switchRegister.addEventListener("click", toggleRegister);

function toggleLogin() {
    register.classList.remove("active");
    login.classList.add("active");
}

function toggleRegister() {
    login.classList.remove("active");
    register.classList.add("active");
}

