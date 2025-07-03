document.getElementById('speakBtn').addEventListener('click', () => {
    const textToSay = document.getElementById('textToSay').value;
    if (textToSay) {
        const audio = new Audio(`/say?say=${encodeURIComponent(textToSay)}`);
        audio.play();
    }
});
