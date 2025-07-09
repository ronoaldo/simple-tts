const shortcuts = [
    "Sim",
    "Não",
    "Obrigado",
    "Por favor",
    "Estou bem!",
    "Estou com sede",
    "Estou com dor"
];

const shortcutsContainer = document.getElementById('shortcuts');
const speakBtn = document.getElementById('speakBtn');
const spinner = speakBtn.querySelector('.spinner-border');
const voiceSelect = document.getElementById('voice');

function playAudio(text) {
    speakBtn.disabled = true;
    spinner.classList.remove('d-none');

    const voice = voiceSelect.value;
    const audio = new Audio(`/say?say=${encodeURIComponent(text)}&voice=${encodeURIComponent(voice)}`);
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

async function loadVoices() {
    try {
        const response = await fetch('/voices');
        const voices = await response.json();
        voices.forEach(voice => {
            const option = document.createElement('option');
            option.value = voice.name;
            option.textContent = voice.name;
            if (voice.name.includes('Fenrir')) {
                option.selected = true;
            }
            voiceSelect.appendChild(option);
        });
    } catch (error) {
        console.error('Failed to load voices:', error);
        alert('Failed to load voices.');
    }
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

loadVoices();
