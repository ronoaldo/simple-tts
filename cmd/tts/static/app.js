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
        let voices = await response.json();

        // Sort voices by gender (male first)
        voices.sort((a, b) => {
            const genderOrder = {
                1: 1, // SSML_GENDER_MALE
                2: 2, // SSML_GENDER_FEMALE
                0: 3, // SSML_GENDER_UNSPECIFIED
                3: 4, // SSML_GENDER_NEUTRAL
            };
            return genderOrder[a.ssml_gender] - genderOrder[b.ssml_gender];
        });

        const maleOptgroup = document.createElement('optgroup');
        maleOptgroup.label = 'Masculino';
        const femaleOptgroup = document.createElement('optgroup');
        femaleOptgroup.label = 'Feminino';
        const neutralOptgroup = document.createElement('optgroup');
        neutralOptgroup.label = 'Neutro';

        voices.forEach(voice => {
            const option = document.createElement('option');
            option.value = voice.name;
            let genderLabel = '';
            switch (voice.ssml_gender) {
                case 1: // SSML_GENDER_MALE
                    genderLabel = '(Masculino)';
                    break;
                case 2: // SSML_GENDER_FEMALE
                    genderLabel = '(Feminino)';
                    break;
                default:
                    genderLabel = '(Neutro)';
            }
            option.textContent = `${voice.name.replace('pt-BR-Chirp3-HD-', '')} ${genderLabel}`.trim();
            if (voice.name.includes('Fenrir')) {
                option.selected = true;
            }

            switch (voice.ssml_gender) {
                case 1: // SSML_GENDER_MALE
                    maleOptgroup.appendChild(option);
                    break;
                case 2: // SSML_GENDER_FEMALE
                    femaleOptgroup.appendChild(option);
                    break;
                default:
                    neutralOptgroup.appendChild(option);
            }
        });

        if (maleOptgroup.children.length > 0) {
            voiceSelect.appendChild(maleOptgroup);
        }
        if (femaleOptgroup.children.length > 0) {
            voiceSelect.appendChild(femaleOptgroup);
        }
        if (neutralOptgroup.children.length > 0) {
            voiceSelect.appendChild(neutralOptgroup);
        }

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
