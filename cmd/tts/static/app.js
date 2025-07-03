const shortcuts = [
    "Obrigado",
    "Por favor",
    "Estou com sede",
    "Estou com dor"
];

const shortcutsContainer = document.getElementById('shortcuts');
const speakBtn = document.getElementById('speakBtn');
const spinner = speakBtn.querySelector('.spinner-border');

function playAudio(text) {
    speakBtn.disabled = true;
    spinner.classList.remove('d-none');

    const audio = new Audio(`/say?say=${encodeURIComponent(text)}`);
    audio.addEventListener('canplaythrough', () => {
        audio.play();
    });
    audio.addEventListener('ended', () => {
        speakBtn.disabled = false;
        spinner.classList.add('d-none');
    });
    audio.addEventListener('error', () => {
        speakBtn.disabled = false;
        spinner.classList.add('d-none');
        alert('Failed to play audio.');
    });
}

shortcuts.forEach(shortcut => {
    const button = document.createElement('button');
    button.textContent = shortcut;
    button.classList.add('btn', 'btn-secondary', 'm-1');
    button.addEventListener('click', () => {
        playAudio(shortcut);
    });
    shortcutsContainer.appendChild(button);
});

document.getElementById('speakBtn').addEventListener('click', () => {
    const textToSay = document.getElementById('textToSay').value;
    if (textToSay) {
        playAudio(textToSay);
    }
});
