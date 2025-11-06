## BETS API
Test jobs of bets api.
```
Github:  https://github.com/kkravchenko/bets_api
``` 


### Run
Clone project and run ` make dev` command.
Go to url: http://localhost:8981/

#### All comands:

| Command          | Description                |
|------------------|----------------------------|
| `make dev`       | Запуск в dev-режиме        |
| `make build`     | Сборка в `/build/bets-api` |
| `make docker-up` | Поднять Docker-контейнер   |
| `make clean`     | Удалить папку `build`      |
| `make help`      | Показать эту справку       |


### END POINTS
POST: /bets - Add bet

```
fetch('http://localhost:8981/bets', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    user_id: "player777",
    amount: 250.50,
    crash_point: 2.75
  })
})
.then(response => {
  if (!response.ok) throw new Error('Ошибка сети');
  return response.json();
})
.then(data => {
  console.log('Ставка создана!', data);
  alert(`ID ставки: ${data.id}`);
})
.catch(err => {
  console.error('Ошибка:', err);
  alert('Не удалось создать ставку');
});
```

GET: /bets - Get bets

```
fetch('http://localhost:8981/bets')
  .then(res => res.json())
  .then(data => {
    console.log('Все ставки:', data);
  });
```

GET: /bets/{id} - Get bet by id

```
const betId = "ff635bba-d17d-4d46-8702-390efe3408e5";

fetch(`http://localhost:8981/bets/${betId}`)
  .then(res => {
    if (!res.ok) throw new Error("Ставка не найдена");
    return res.json();
  })
  .then(bet => {
    console.log("Ставка найдена! 🎯", bet);
  })
  .catch(err => {
    console.error(err);
  });
```

GET: /health - Check health

```
fetch('http://localhost:8981/health')
  .then(res => res.text())
  .then(status => {
    console.log(status);
  })
  .catch(err => {
    console.error(err);
  });
```