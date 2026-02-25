// ────────────────────────────────────────
// CONFIG
// ────────────────────────────────────────
const API = '';

// session_id : UUID stocké en localStorage pour identifier le panier
// Si l'utilisateur revient, son panier est retrouvé via ce même UUID
function getSessionID() {
    let sid = localStorage.getItem('session_id');
    if (!sid) {
    // Génère un UUID v4 simple sans dépendance externe
    sid = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, c => {
        const r = Math.random() * 16 | 0;
        return (c === 'x' ? r : (r & 0x3 | 0x8)).toString(16);
    });
    localStorage.setItem('session_id', sid);
    }
    return sid;
}

const SESSION = getSessionID();

// ────────────────────────────────────────
// UTILITAIRES
// ────────────────────────────────────────

// Affiche un message court en bas d'écran pendant 2 secondes
function showToast(msg) {
    const t = document.getElementById('toast');
    t.textContent = msg;
    t.classList.add('show');
    setTimeout(() => t.classList.remove('show'), 2000);
}

// Met à jour le badge dans le header
function updateCartCount(items) {
    const total = items.reduce((acc, i) => acc + i.quantity, 0);
    document.getElementById('cart-count').textContent = total;
}

// ────────────────────────────────────────
// CATALOGUE - Chargement & rendu
// ────────────────────────────────────────
async function loadGames() {
    try {
    // GET /games → tableau de Game
    const res   = await fetch(`${API}/games`);
    const games = await res.json();

    document.getElementById('loader').style.display = 'none';
    const grid = document.getElementById('games-grid');

    games.forEach(game => {
        // Crée une carte par jeu
        const card = document.createElement('div');
        card.className = 'game-card';
        card.innerHTML = `
        <img src="${game.image_url}" alt="${game.title}" loading="lazy" />
        <div class="card-body">
            <div class="card-title">${game.title}</div>
            <div class="card-desc">${game.description}</div>
            <div class="card-footer">
            <span class="price">${game.price.toFixed(2)} €</span>
            <button class="add-btn" data-id="${game.id}" ${game.stock === 0 ? 'disabled' : ''}>
                ${game.stock === 0 ? 'Rupture' : 'Ajouter'}
            </button>
            </div>
        </div>
        `;
        grid.appendChild(card);
    });

    // Délégation d'événement : un seul listener pour tous les boutons "Ajouter"
    // Plus efficace que d'attacher un listener à chaque carte individuellement
    grid.addEventListener('click', async (e) => {
        if (!e.target.classList.contains('add-btn')) return;
        const gameID = parseInt(e.target.dataset.id);
        await addToCart(gameID);
    });

    } catch (err) {
    document.getElementById('loader').textContent = 'Impossible de charger les jeux.';
    console.error(err);
    }
}

// ────────────────────────────────────────
// PANIER - Opérations
// ────────────────────────────────────────

// Ajoute un jeu au panier via POST
async function addToCart(gameID) {
    await fetch(`${API}/cart/${SESSION}/add`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    // Body : game_id + quantity (1 par défaut)
    body: JSON.stringify({ game_id: gameID, quantity: 1 })
    });
    showToast('Ajouté au panier !');
    await loadCart(); // Recharge le panier pour mettre à jour l'affichage
}

// Supprime un jeu du panier via DELETE
async function removeFromCart(gameID) {
    await fetch(`${API}/cart/${SESSION}/remove/${gameID}`, { method: 'DELETE' });
    await loadCart();
}

// Charge et affiche le panier courant
async function loadCart() {
    const res  = await fetch(`${API}/cart/${SESSION}`);
    const cart = await res.json();

    updateCartCount(cart.items);

    const container = document.getElementById('cart-items');
    const emptyMsg  = document.getElementById('cart-empty');

    // Supprime les anciennes lignes (mais pas le message vide)
    container.querySelectorAll('.cart-item').forEach(el => el.remove());

    if (!cart.items || cart.items.length === 0) {
    emptyMsg.style.display = 'block';
    document.getElementById('cart-total').textContent = '0.00 €';
    return;
    }

    emptyMsg.style.display = 'none';

    // Affiche chaque article
    cart.items.forEach(item => {
    const div = document.createElement('div');
    div.className = 'cart-item';
    div.innerHTML = `
        <img src="${item.image_url}" alt="${item.title}" />
        <div class="cart-item-info">
        <div class="title">${item.title}</div>
        <div class="qty">x${item.quantity}</div>
        </div>
        <span class="cart-item-price">${item.subtotal.toFixed(2)} €</span>
        <button class="remove-btn" data-id="${item.game_id}">✕</button>
    `;
    container.appendChild(div);
    });

    document.getElementById('cart-total').textContent = cart.total.toFixed(2) + ' €';

    // Listener pour les boutons de suppression
    container.querySelectorAll('.remove-btn').forEach(btn => {
    btn.addEventListener('click', () => removeFromCart(parseInt(btn.dataset.id)));
    });
}

// ────────────────────────────────────────
// PANNEAU PANIER - Ouverture / fermeture
// ────────────────────────────────────────
function openCart() {
    document.getElementById('cart-overlay').classList.add('open');
    document.getElementById('cart-panel').classList.add('open');
    loadCart();
}

function closeCart() {
    document.getElementById('cart-overlay').classList.remove('open');
    document.getElementById('cart-panel').classList.remove('open');
}

document.getElementById('cart-btn').addEventListener('click', openCart);
document.getElementById('close-cart').addEventListener('click', closeCart);
// Ferme le panier en cliquant sur l'overlay (en dehors du panel)
document.getElementById('cart-overlay').addEventListener('click', closeCart);

// ────────────────────────────────────────
// INIT
// ────────────────────────────────────────
loadGames();
loadCart(); // Charge le panier au démarrage pour mettre à jour le badge