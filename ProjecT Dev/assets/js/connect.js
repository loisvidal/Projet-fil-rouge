document.addEventListener("DOMContentLoaded", () => {
    const connectContainer = document.querySelector('.Connect');
    const body = document.querySelector('body');
    const registerForm = document.querySelector('.register form');
    
    // Boutons du Header
    const btnHeaderLogin = document.querySelector('.isConnect_login');
    const btnHeaderRegister = document.querySelector('.isConnect_register');


    // 1. Injection des liens de switch internes
    const loginSwitchHTML = `
        <div class="switch-text">
            Pas encore de compte ? <a class="switch-link" id="go-to-register">S'inscrire</a>
        </div>
    `;
    document.querySelector('.login').insertAdjacentHTML('beforeend', loginSwitchHTML);

    const registerSwitchHTML = `
        <div class="switch-text">
            Déjà un compte ? <a class="switch-link" id="go-to-login">Se connecter</a>
        </div>
    `;
    document.querySelector('.register').insertAdjacentHTML('beforeend', registerSwitchHTML);

    const btnGoRegister = document.getElementById('go-to-register');
    const btnGoLogin = document.getElementById('go-to-login');

    // --- LOGIQUE D'OUVERTURE / FERMETURE ---

    function openModal(isRegister = false) {
        connectContainer.classList.add('show');
        body.classList.add('modal-open'); // Active le fond noir flouté
        
        // Si on a cliqué sur Register, on glisse directement sur la bonne slide
        if (isRegister) {
            connectContainer.classList.add('slide-active');
        } else {
            connectContainer.classList.remove('slide-active');
        }
    }

    function closeModal() {
        connectContainer.classList.remove('show');
        body.classList.remove('modal-open');
    }

    // --- ÉCOUTEURS D'ÉVÉNEMENTS ---

    // Clic sur "Login" dans le Header
    if (btnHeaderLogin) {
        btnHeaderLogin.addEventListener('click', (e) => {
            console.log("titi");
            e.preventDefault(); // Bloque le comportement par défaut au cas où tu as mis un lien
            openModal(false);
        });
    }

    // Clic sur "Register" dans le Header
    if (btnHeaderRegister) {
        btnHeaderRegister.addEventListener('click', (e) => {
            console.log("toto");
            e.preventDefault();
            openModal(true);
        });
    }

    // Clic sur les liens internes pour switcher
    btnGoRegister.addEventListener('click', () => {
        connectContainer.classList.add('slide-active');
    });

    btnGoLogin.addEventListener('click', () => {
        connectContainer.classList.remove('slide-active');
    });

    // Retourne au login lors de la soumission du Register (si on ne change pas de page)
    registerForm.addEventListener('submit', () => {
        connectContainer.classList.remove('slide-active');
    });

    // Bonus : Fermer la modale si on appuie sur la touche "Échap"
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape' && connectContainer.classList.contains('show')) {
            closeModal();
        }
    });
});