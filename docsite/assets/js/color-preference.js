const lsKeyColorPreference = 'color-preference'

const getColorPreference = () => {
    let lastUsedColorPreference = localStorage.getItem(lsKeyColorPreference)
    if (lastUsedColorPreference !== null)
        return lastUsedColorPreference
    else
        return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

let colorPreference = getColorPreference()
localStorage.setItem(lsKeyColorPreference, colorPreference)
document.firstElementChild.setAttribute('data-color', colorPreference)

const colorPreferenceButton = document.querySelector('#content-wrapper > footer > div:last-child > button:last-child')
if (colorPreferenceButton) {
    colorPreferenceButton.addEventListener('click', function () {
        if (colorPreference == 'dark') {
            colorPreference = 'light'
        } else {
            colorPreference = 'dark'
        }
        setColorPreference()
    });
}

const setColorPreference = () => {
    localStorage.setItem(lsKeyColorPreference, colorPreference)
    document.firstElementChild.setAttribute('data-color', colorPreference)
}

window.matchMedia('(prefers-color-scheme: dark)')
    .addEventListener('change', ({matches: isDark}) => {
        colorPreference = isDark ? 'dark' : 'light'
        setColorPreference()
    })