var html = document.querySelector('html')
var title = document.getElementById('title')
var actionsBtn = document.querySelector('#todo__action__bar > button')
var colorsBar = document.querySelector('.todo__action__color')
var input = document.getElementById('todo-write')
var form = document.getElementById('form')
var todoList = document.getElementById('todo-list')
var doneList = document.getElementById('done-list')
var listId = location.hash
var delay = 0
var timeout
var sock
var wsp = location.protocol === 'https:' ? 'wss' : 'ws'
;
autocomplete({
    input: input,
    fetch: function(text, update) {
        text = text.toLowerCase()
        var suggestions = []
        var doneItems = doneList.querySelectorAll('.todo__text')
        var todoItems = todoList.querySelectorAll('.todo__text')

        doneItems.forEach(el => {
            const value = el.innerText
            if (value.toLowerCase().startsWith(text))
                suggestions.push({ value, label: value })
        })

        todoItems.forEach(el => {
            const value = el.innerText
            if (value.toLowerCase().startsWith(text))
                suggestions.push({ value, label: value })
        })

        update(suggestions)
    },
    onSelect: function(item) {
        input.value = item.label
        onAddItem()
    }
})
;

function startWebsocket() {
    sock = new WebSocket(`${wsp}://${location.host}/ws`);

    sock.onclose = function (event) {
        document.body.style.opacity = .5
        input.setAttribute('disabled', true)

        console.log('connection closed')

        sock = null
        delay = (delay > 10*1000) ? delay : (delay + 200)
        timeout = setTimeout(startWebsocket, delay)
    }

    sock.onmessage = function (event) {
        try {
            var data = [].concat(JSON.parse(event.data));
            data.forEach(item => {

                switch (item.method) {
                    case 'Ping':
                        console.log('connection openned')
                        input.removeAttribute('disabled')
                        document.body.style.opacity = 1
                        break;

                    case 'DeleteItem':
                        var exists = document.getElementById(`item-${item.itemId}`)
                        if (exists) {
                            document.getElementById(`item-${item.itemId}`).remove()
                        }
                        break;

                    case 'AddItem':
                        var exists = document.getElementById(`item-${item.itemId}`)
                        if (exists) {
                            document.getElementById(`item-${item.itemId}`).remove()
                        }
                        var node = document.createElement('div')
                        node.innerHTML = itemTpl(item)
                        setTimeout(registerChangeCallbackFn(node.firstElementChild, item), 1);
                        setTimeout(registerDeleteCallbackFn(node.firstElementChild, item), 1);
                        (item.isChecked ? doneList : todoList).prepend(node.firstElementChild)
                        break;

                    default:
                        console.log('unhandled', item)
                }
            })
        } catch (e) {
            console.error(e)
        }
    }

    sock.onopen = function () {
        console.log('connection opened')
        sock.send(JSON.stringify({
            method: 'GetItems',
            listId
        }))
    }
}

if (!location.hash) {
    location.hash = listId = `#${uuid()}:New List`;
    location.reload()
} else {
    hashAndTitle = location.hash.split(':')
    title.innerText = decodeURIComponent(hashAndTitle[1] ? hashAndTitle[1] : hashAndTitle[0])
    listId = hashAndTitle[0]
    startWebsocket()
}
window.addEventListener('offline', function(e) {
    console.log('offline')
    clearTimeout(timeout)
    if (sock) sock.onclose()
});
window.addEventListener('online', function(e) {
    console.log('online');
    clearTimeout(timeout)
    if (sock) sock.onclose()
    startWebsocket()
});

document.querySelectorAll('.todo__action__color button').forEach(el => {
    el.addEventListener('click', evt => onChangeTheme(evt.target.style.backgroundColor))
})

html.addEventListener('click', evt => {
    if (evt.target.tagName.toLowerCase() != 'button')
        colorsBar.style.display = 'none', {capture:true}
})
actionsBtn.addEventListener('click', _ => {
    const visible = colorsBar.style.display == 'block';
    colorsBar.style.display = visible ? 'none' : 'block'
    actionsBtn.innerText = visible ? '☰' : '—'
})

title.addEventListener('click', _ => {
    title.contentEditable = true;
    title.focus()
})

title.addEventListener('keypress', evt => {
    if (event.key === "Enter") {
        event.preventDefault();
        onChangeTitle(title.innerText)
    }
})

function reconnectWebsocket() {
    if (sock) sock.onclose()
    startWebsocket()
}

function registerChangeCallbackFn(node, item) {
    return _ => node.querySelector('input').addEventListener('change', event => {
        return onUpdateItem(event, item.itemId, event.target.checked)
    })
}

function registerDeleteCallbackFn(node, item) {
    return _ => node.querySelector('.todo__delete').addEventListener('click', event => {
        node.classList.add('todo_removed')
        return onDeleteItem(event, item.itemId, event.target.checked)
    })
}

function onAddItem() {
    event.preventDefault()
    var text = ''+input.value
    if (!text.trim()) return;
    form.reset()
    sock.send(JSON.stringify({
        method: 'AddItem',
        isChecked: false,
        listId,
        text,
    }))
    return false
}

function onUpdateItem(event, itemId, isChecked) {
    event.preventDefault()
    setTimeout(_ => {
        sock.send(JSON.stringify({
            method: 'UpdateItem',
            isChecked: !!isChecked,
            itemId,
            listId,
        }))
    }, 400)
    return false
}

function onDeleteItem(event, itemId) {
    event.preventDefault()
    setTimeout(_ => {
        sock.send(JSON.stringify({
            method: 'DeleteItem',
            itemId,
            listId,
        }))
    }, 400)
    return false
}

function onChangeTitle(text) {
    title.innerText = text
    title.contentEditable = false
    location.hash = listId + ':' + text
    localStorage.setItem(listId, text)
    favicon(html.style.backgroundColor, text.slice(0, 1))
}

function onChangeTheme(color) {
    html.style.backgroundColor = color
    localStorage.setItem('theme', color)
}

function uuid(a,b){for(b=a='';a++<8;b+=a*51&52?(a^15?8^Math.random()*(a^20?16:4):4).toString(16):'-');return b}

function itemTpl(item) {
    return `
    <label class="todo" id="item-${item.itemId}">
        <input class="todo__state" type="checkbox" id="check-${item.itemId}" ${item.isChecked ? 'checked' : ''} />

        <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 200 25"
            class="todo__icon">
            <use xlink:href="#todo__line" class="todo__line"></use>
            <use xlink:href="#todo__box" class="todo__box"></use>
            <use xlink:href="#todo__check" class="todo__check"></use>
            <use xlink:href="#todo__circle" class="todo__circle"></use>
        </svg>

        <div class="todo__text">${item.text}</div>

        <div class="todo__delete">⌦</div>
    </label>`;
}

function favicon(color, letter) {
    var canvas = document.createElement('canvas');
    canvas.width = 16;
    canvas.height = 16;
    var ctx = canvas.getContext('2d');
    var img = new Image();
    // img.src = '/favicon.ico';

    ctx.fillStyle = color;
    ctx.fillRect(0, 0, canvas.width, canvas.height);

    ctx.fillStyle = '#FFFFFF';
    ctx.font = 'bold 10px sans-serif';
    ctx.fillText(letter, 4, 12);

    var link = document.createElement('link');
    link.type = 'image/x-icon';
    link.rel = 'shortcut icon';
    link.href = canvas.toDataURL("image/x-icon");

    document.querySelectorAll('[rel="shortcut icon"]').forEach(el => el.remove())
    document.getElementsByTagName('head')[0].appendChild(link);
}

/****************************/

var color = localStorage.getItem('theme')
if (color) html.style.backgroundColor = color;

var customTitle = localStorage.getItem(listId)
if (customTitle) onChangeTitle(customTitle);
else title.click()
