fetch('/api/items')
    .then(res => res.json())
    .then(data => {
        const list = document.getElementById('items');
        list.innerHTML = '';
        data.forEach(item => {
            const li = document.createElement('li');
            li.textContent = `${item.id}: ${item.name}`;
            list.appendChild(li);
        });
    })
    .catch(err => {
        document.getElementById('items').textContent = 'Ошибка загрузки данных';
        console.error(err);
    });