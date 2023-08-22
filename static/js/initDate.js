(() => {
    const inputGiorno = document.getElementById("giornoNascitaInput");
    const inputMese = document.getElementById("meseNascitaInput");
    const inputAnno = document.getElementById("annoNascitaInput");
    const mappaMese = {
        2: "Febbraio",
        3: "Marzo",
        4: "Aprile",
        5: "Maggio",
        6: "Giugno",
        7: "Luglio",
        8: "Agosto",
        9: "Settembre",
        10: "Ottobre",
        11: "Novembre",
        12: "Dicembre"
    };

    // Add days of the month to select form
    Array(32).fill(0).map((_, i) => i).slice(2).forEach(day => {
        const newOption = document.createElement("option");
        newOption.text = day.toString();
        inputGiorno.add(newOption);
    });

    // Add months of the year to select form
    Array(13).fill(0).map((_, i) => i).slice(2).forEach(month => {
        const newOption = document.createElement("option");
        newOption.value = month.toString();
        newOption.label = mappaMese[month];
        newOption.text = mappaMese[month];
        inputMese.add(newOption);
    });

    // Add years to select form
    const startYear = 1901;
    const endYear = new Date().getFullYear();
    [...Array(endYear - startYear + 1).keys()].map(i => i + startYear).forEach(year => {
        const newOption = document.createElement("option");
        newOption.text = year.toString();
        inputAnno.add(newOption);
    });
    inputAnno.value = "1985";
})();