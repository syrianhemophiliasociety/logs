<div align="center">
  <a href="https://syrianhemophiliasociety.com" target="_blank"><img src="https://syrianhemophiliasociety.com/assets/android-chrome-512x512.png" width="150" /></a>

  <h1>SyrianHemophiliaSocietyLogs</h1>
  <p>
    <strong>A patient care follow up platform for Syrian Hemophilia Society</strong>
  </p>
  <p>
    <a href="https://goreportcard.com/report/github.com/mbaraa/shs"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/mbaraa/shs"/></a>
    <a href="https://github.com/mbaraa/shs/actions/workflows/rex-deploy.yml"><img alt="rex-deployment" src="https://github.com/mbaraa/shs/actions/workflows/rex-deploy.yml/badge.svg"/></a>
  </p>
</div>

## About

**SyrianHemophiliaSocietyLogs** is something idk.

_Note: this is a fling side-project that could die anytime so don't get your hopes up._

## Contributing

IDK, it would be really nice of you to contribute, check the poorly written [CONTRIBUTING.md](/CONTRIBUTING.md) for more info.

## Run locally

1. Clone the repo.

```bash
git clone https://github.com/mbaraa/shs
```

2. Create the docker environment file

```bash
cp .env.example .env.docker
```

3. Run it with docker compose.

```bash
docker compose up -f docker-compose-all.yml
```

3. Visit http://localhost:23103
4. Don't ask why I chose this weird port.

---

A [DankStuff <img height="16" width="16" src="https://dankstuff.net/assets/favicon.ico" />](https://dankstuff.net) product!

Made with 🧉 by [Baraa Al-Masri](https://mbaraa.com)
