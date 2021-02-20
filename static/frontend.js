(function(window){

    function uuid(a,b){for(b=a='';a++<8;b+=a*51&52?(a^15?8^Math.random()*(a^20?16:4):4).toString(16):'-');return b}
    var $ = document.querySelector.bind(document);
    var $$ = document.querySelectorAll.bind(document);
    
    var UI = {
        // current list id
        listId: location.hash,

        ws: null,
    
        // UI elements
        html: $('html'),
        head: $('head'),
        form: $('#form'),
        title: $('#title'),
        input: $('#todo-write'),
        todoList: $('#todo-list'),
        doneList: $('#done-list'),
        colorsBar: $('.todo__action__color'),
        actionsBtn: $('#todo__action__bar > button'),
    
        // temporary containers
        shadowTodo: document.createElement('div'),
        shadowDone: document.createElement('div'),

        init(ws) {
            this.ws = ws;

            // Actions bar event handlers

            // Toggle the menu when press the nemu button
            this.actionsBtn.addEventListener('click', _ => {
                const visible = this.colorsBar.style.display == 'block';
                this.colorsBar.style.display = visible ? 'none' : 'block'
                this.actionsBtn.innerText = visible ? '☰' : '—'
            })
            // Close the menu when press anywhere else
            this.html.addEventListener('click', evt => {
                if (evt.target.tagName.toLowerCase() != 'button')
                    this.colorsBar.style.display = 'none', {capture:true}
            })
            // Change theme buttons
            $$('.todo__action__color button').forEach(el => {
                el.addEventListener('click', evt => this.onChangeTheme(evt.target.style.backgroundColor))
            })
            
            // Title event handlers

            // Make the title editable
            this.title.addEventListener('click', _ => {
                this.title.contentEditable = true;
                this.title.focus()
            })
            // Update the title on ENTER
            this.title.addEventListener('keypress', evt => {
                if (evt.key === "Enter") {
                    evt.preventDefault();
                    this.onChangeTitle(this.title.innerText)
                }
            })
    
            // Form & Autocomplete

            // Add a new item on submit
            this.form.addEventListener('submit', evt => {
                evt.preventDefault()
                this.ws.onAddItem(this.input.value.trim())
                return false;
            })
            // Init autocompleter
            this.autocomplete()
        },
    
        // Enable interface
        enable() {
            document.body.style.opacity = 1
            this.input.removeAttribute('disabled')
        },

        // Disable interface
        disable() {
            document.body.style.opacity = .5
            this.input.setAttribute('disabled', true)
        },

        // Add a Update listener on this node
        registerChangeCallbackFn(node, item) {
            return _ => node.querySelector('input').addEventListener('change', event => {
                return this.ws.onUpdateItem(event, item.itemId, event.target.checked)
            })
        },
    
        // Add a Delete listener on this node
        registerDeleteCallbackFn(node, item) {
            return _ => node.querySelector('.todo__delete').addEventListener('click', event => {
                node.classList.add('todo_removed')
                return this.ws.onDeleteItem(event, item.itemId, event.target.checked)
            })
        },
    
        // Change title callback
        onChangeTitle(text) {
            this.title.innerText = text
            this.title.contentEditable = false
            location.hash = this.listId + ':' + text
            // store the title in localStorage
            localStorage.setItem(this.listId, text)
            // make a pretty favicon
            this.redrawFavicon(this.html.style.backgroundColor, text.slice(0, 1))
        },
    
        // Change theme callback
        onChangeTheme(color) {
            $$('.themed').forEach(el => el.style.backgroundColor = color);
            localStorage.setItem('theme', color)
        },
    
        // Called on each received item
        onItemReceived(item) {
            switch (item.method) {
                case 'Ping':
                    console.log('connection openned')
                    this.enable()
                    break;
    
                case 'DeleteItem':
                    var exists = $(`#item-${item.itemId}`)
                    if (exists) {
                        $(`#item-${item.itemId}`).remove()
                    }
                    break;
    
                case 'AddItem':
                    var exists = $(`#item-${item.itemId}`)
                    if (exists) {
                        $(`#item-${item.itemId}`).remove()
                    }
                    var node = document.createElement('div')
                    node.innerHTML = this.itemTemplate(item)
                    setTimeout(this.registerChangeCallbackFn(node.firstElementChild, item), 1);
                    setTimeout(this.registerDeleteCallbackFn(node.firstElementChild, item), 1);
                    (item.isChecked ? this.shadowDone : this.shadowTodo).prepend(node.firstElementChild)
                    break;
    
                default:
                    console.warn('unhandled', item)
            }
        },

        commitItemReceived() {
            this.doneList.prepend.apply(this.doneList, this.shadowDone.children)
            this.todoList.prepend.apply(this.todoList, this.shadowTodo.children)
        },

        // form autocompleter
        autocomplete() {
            var self = this;
            return autocomplete({
                input: self.input,
                fetch: function(text, update) {
                    text = text.toLowerCase()
                    var suggestions = []
                    var doneItems = self.doneList.querySelectorAll('.todo__text')
                    var todoItems = self.todoList.querySelectorAll('.todo__text')
        
                    // complete on todo and done items

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
                    self.input.value = item.label
                    self.ws.onAddItem(self.input.value.trim())
                }
            })
        },

        itemTemplate(item) {
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
                <svg class="todo__delete" viewBox="0 0 24 24" width="24" height="24" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round" class="css-i6dzq1">
                    <line x1="18" y1="6" x2="6" y2="18"></line>
                    <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
            </label>`;
        },
    
        redrawFavicon(color, letter) {
            var canvas = document.createElement('canvas');
            canvas.width = 16;
            canvas.height = 16;
            var ctx = canvas.getContext('2d');
    
            ctx.fillStyle = color;
            ctx.fillRect(0, 0, canvas.width, canvas.height);
    
            ctx.fillStyle = '#FFFFFF';
            ctx.font = 'bold 10px sans-serif';
            ctx.fillText(letter, 4, 12);
    
            var link = document.createElement('link');
            link.type = 'image/x-icon';
            link.rel = 'shortcut icon';
            link.href = canvas.toDataURL("image/x-icon");
    
            $$('[rel="shortcut icon"]').forEach(el => el.remove())
            this.head.appendChild(link);
        },
    };
    
    var WS = {
        ui: null,

        // websocket
        sock: undefined,
        wsp: location.protocol === 'https:' ? 'wss' : 'ws',
    
        // WS delay to next reconnection
        timeout: undefined,
        delay: 0,
    
        // init
        init(ui) {
            this.ui = ui;
            window.addEventListener('offline', this.offline)
            window.addEventListener('online', this.online)
            this.connect()
        },
    
        // ws handlers
        onClose() {
            console.log('connection closed')
            this.ui.disable()
    
            this.sock = null;
            this.delay = (this.delay > 10*1000) ? this.delay : (this.delay + 200)
            this.timeout = setTimeout(this.startWebsocket.bind(this), this.delay)
        },
    
        onMessage(event) {
            try {
                var data = [].concat(JSON.parse(event.data));
                data.forEach(item => this.ui.onItemReceived(item))
                this.ui.commitItemReceived();
            } catch (e) {
                console.error(e)
            }
        },
    
        onOpen() {
            console.log('connection opened')
            this.sock.send(JSON.stringify({
                method: 'GetItems',
                listId: this.ui.listId,
            }))
        },
    
        connect() {
            if (this.sock) this.sock.close()
            this.startWebsocket()
        },

        startWebsocket() {
            this.sock = new WebSocket(`${this.wsp}://${location.host}/ws`);
            this.sock.onmessage = this.onMessage.bind(this);
            this.sock.onclose = this.onClose.bind(this);
            this.sock.onopen = this.onOpen.bind(this);
        },
    
        offline(e) {
            console.log('offline')
            clearTimeout(this.timeout)
            if (this.sock) this.sock.close()
        },
    
        online(e) {
            console.log('online');
            clearTimeout(this.timeout)
            this.connect()
        },
    
        // functional handlers
        onAddItem(text) {
            if (!text.trim()) return;
            this.ui.form.reset()
            this.sock.send(JSON.stringify({
                method: 'AddItem',
                isChecked: false,
                listId: this.ui.listId,
                text,
            }))
            return false
        },
        onUpdateItem(event, itemId, isChecked) {
            event.preventDefault()
            setTimeout(_ => {
                this.sock.send(JSON.stringify({
                    method: 'UpdateItem',
                    isChecked: !!isChecked,
                    listId: this.ui.listId,
                    itemId,
                }))
            }, 400)
            return false
        },
        onDeleteItem(event, itemId) {
            event.preventDefault()
            setTimeout(_ => {
                this.sock.send(JSON.stringify({
                    method: 'DeleteItem',
                    listId: this.ui.listId,
                    itemId,
                }))
            }, 400)
            return false
        },
    }
    
    // UI & transport
    if (!location.hash) {
        // no hash, set a random one and reload
        location.hash = UI.listId = `#${uuid()}:New List`;
        location.reload()

    } else {
        // app bootstrap
        hashAndTitle = location.hash.split(':')
        UI.title.innerText = decodeURIComponent(hashAndTitle[1] ? hashAndTitle[1] : hashAndTitle[0])
        UI.listId = hashAndTitle[0]

        UI.init(WS)
        WS.init(UI)
    }

    // Theme
    var color = localStorage.getItem('theme')
    if (color) UI.onChangeTheme(color)

    // List title
    var customTitle = localStorage.getItem(UI.listId)
    if (customTitle) UI.onChangeTitle(customTitle);
    else title.click()
    
})(window)
    
