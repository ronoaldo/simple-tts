document.getElementById('speakBtn').addEventListener('click', () => {
    const textToSay = document.getElementById('textToSay').value;
    const speakBtn = document.getElementById('speakBtn');
    const spinner = speakBtn.querySelector('.spinner-border');

    if (textToSay) {
        speakBtn.disabled = true;
        spinner.classList.remove('d-none');

        const audio = new Audio(`/say?say=${encodeURIComponent(textToSay)}`);
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
});
