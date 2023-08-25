const copy = () => {
    const element = document.getElementById("res");
    const codiceFiscale = element.innerText;

    navigator.clipboard.writeText(codiceFiscale);
    element.innerText = "Copiato nella clipboard!"

    setTimeout(() => {
        element.innerText = codiceFiscale;
    }, 700);
};