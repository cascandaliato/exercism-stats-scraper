# Exercism Stats Scraper · [![Netlify](https://img.shields.io/netlify/dd4af605-f9cf-4b1a-9bcc-d2a35eaa66dd?logo=netlify&style=for-the-badge)](https://app.netlify.com/sites/exercism-stats-scraper/deploys) [![Exercism statistics](https://img.shields.io/badge/dynamic/json?color=009caa&label=casca%27s%20solutions&query=total&url=https%3A%2F%2Fexercism-stats-scraper.netlify.app%2Fapi%2Fsolutions%3Fuser%3Dcasca&logo=exercism&logoColor=white&style=for-the-badge)](https://exercism.io/profiles/casca)
A super simple scraper of [Exercism](http://exercism.io/) user profiles that can be used to create badges with [Shields.io](https://shields.io/) (scroll down for [instructions](#create-a-badge)).

- Right now it retrieves only the number of published solutions from the user profile page `https://exercism.io/profiles/<user>`.
- It consists of a single endpoint published at `https://exercism-stats-scraper.netlify.app/api/solutions`.
- The endpoint accepts `GET` requests and requires a `user` param which must be a valid Exercism user, e.g. https://exercism-stats-scraper.netlify.app/api/solutions?user=casca.
- It returns a JSON object with the `total` number of published solutions:
    ```
    {
      "total": number
    }
    ```
  or, if an error occurred:
    ```
    {
      "error": string
    }
    ```

### Create a Badge
1. Head to [Shields.io](https://shields.io/).
1. Fill the form for the `Dynamic` badge:
    - data type: `json`
    - data url: `https://exercism-stats-scraper.netlify.app/api/solutions?query=<user>` (replace `<user>` with your Exercism user)
    - query: `total`
    - specify the badge label (e.g. `solutions`)
1. Hit the Make Badge button.
1. Use the resulting image URL to publish your badge.
1. Optionally, you can add the Exercism logo to your badge by appending `&logo=exercism&logoColor=white` to the image URL.
