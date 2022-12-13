// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.

// ACE editor
var editor;

// Spinner
var spinner;

function formOnload() {
    // Setup the ACE editor.
    editor = ace.edit("editor");
    editor.getSession().setUseWorker(false);

    // Form elements: Editor Options.
    const settingsInvisibleChars =
        document.getElementById("flexSwitchInvisibleChars");
    settingsInvisibleChars.onclick = () => {
        editor.setOption("showInvisibles", settingsInvisibleChars.checked);
        setCookie( "flexSwitchInvisibleChars", settingsInvisibleChars.checked, 360);
    };

    const settingsDarkEditor = document.getElementById("flexSwitchDarkEditor");
    settingsDarkEditor.onclick = () => {
        if (settingsDarkEditor.checked) {
            editor.setTheme("ace/theme/monokai");
        } else {
            editor.setTheme("ace/theme/cloud9_day");
        }
        setCookie("flexSwitchDarkEditor", settingsDarkEditor.checked, 360);
    };

    // Load settings from cookies (if present) and set initial mode
    // based on the state of the settings form.
    settingsInvisibleChars.checked =
        (hasCookie("flexSwitchInvisibleChars") &&
        getCookie("flexSwitchInvisibleChars") == "true");
    settingsInvisibleChars.onclick();

    settingsDarkEditor.checked =
        (hasCookie("flexSwitchDarkEditor") &&
         getCookie("flexSwitchDarkEditor") == "true");
    settingsDarkEditor.onclick();

    // Set editor language based on language selected in form.
    const languageSet = document.getElementById("languageSelect");
    languageSet.onclick = () => {
        SetACELang(languageSet);
    };
    // Set initial value for the language.
    languageSet.onclick();

    // Spinner
    spinner = document.getElementById("pleasewait");

}

function lint() {
    spinner.style.visibility = "visible";

    var xhttp = new XMLHttpRequest();
    xhttp.open("POST", "{{.LintPath}}", true);

    xhttp.onreadystatechange = function () {
        // Spinner off
        spinner.style.visibility = "hidden";

        if (this.readyState == 4) {
            if (this.status == 200) {
                var res = JSON.parse(this.responseText);
                // Update editor text if code reformatted.
                if (res.Reformatted == true) {
                    editor.setValue(res.ReformattedText);
                }

                if (res.Pass == true) {
                    eid = "results_ok";
                    msg = "No errors found!";
                } else {
                    eid = "results_bad";
                    msg = res.ErrorMessages.join("<br>");
                }
            } else {
                eid = "results_bad";
                msg = "Request failed: " + this.responseText;
            }

            document.getElementById("results_ok_div").style.display = "none";
            document.getElementById("results_bad_div").style.display = "none";
            document.getElementById(eid + "_div").style.display = "block";
            document.getElementById(eid).innerHTML = msg;
        }
    };

    // Send
    let programText = encodeURIComponent(editor.getValue());
    let lang = document.getElementById("languageSelect");
    let req = JSON.stringify({"lang": lang.value, "text": programText});

    xhttp.setRequestHeader("Content-type", "application/json");
    xhttp.send(req);
}

// SetACELang sets the language used by the ACE editor.
function SetACELang(langobj) {
    // Ugly hack: ACE considers C and C++ a single language: c_cpp
    let chosenLang = langobj.options[langobj.selectedIndex].value;
    if (chosenLang == "c" || chosenLang == "cpp") {
        chosenLang = "c_cpp";
    }
    editor.getSession().setMode("ace/mode/" + chosenLang);
}
// Checks the existence of a given cookie.
function hasCookie(cname) {
    let name = cname + "=";
    let ca = document.cookie.split(";");
    for (let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) == " ") {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return true;
        }
    }
    return false;
}

// Retrieve the value of a cookie.
function getCookie(cname) {
    let name = cname + "=";
    let ca = document.cookie.split(";");
    for (let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) == " ") {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}

// Set a cookie with a givcen value and expiration (in days).
function setCookie(cname, cvalue, exdays) {
    const d = new Date();
    d.setTime(d.getTime() + (exdays * 24 * 60 * 60 * 1000));
    let expires = "expires=" + d.toUTCString();
    document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/";
}

// vim: ts=4:sw=4:expandtab:smarttab
