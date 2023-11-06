(function(window){

    function uuid(a,b){for(b=a='';a++<8;b+=a*51&52?(a^15?8^Math.random()*(a^20?16:4):4).toString(16):'-');return b}
    var $ = document.querySelector.bind(document);
    var $$ = document.querySelectorAll.bind(document);
    
    var IconDelete = `<svg class="todo__delete" viewBox="0 0 24 24" width="24" height="24" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round" class="css-i6dzq1"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>`;
    var IconCreate = `<svg viewBox="0 0 512 512" xml:space="preserve" height="1.8em"><g><path style="fill:#F4F6F9;" d="M173.234,262.897h75.869v75.869c0,3.809,3.085,6.897,6.897,6.897c3.812,0,6.897-3.088,6.897-6.897   v-75.869h75.869c3.812,0,6.897-3.088,6.897-6.897c0-3.809-3.085-6.897-6.897-6.897h-75.869v-75.869   c0-3.809-3.085-6.897-6.897-6.897c-3.812,0-6.897,3.088-6.897,6.897v75.869h-75.869c-3.812,0-6.897,3.088-6.897,6.897   C166.337,259.809,169.422,262.897,173.234,262.897z"></path><path style="fill:#F4F6F9;" d="M414.635,90.468H97.365c-3.812,0-6.897,3.088-6.897,6.897v317.269c0,3.809,3.085,6.897,6.897,6.897   h317.269c3.812,0,6.897-3.088,6.897-6.897V97.365C421.532,93.556,418.447,90.468,414.635,90.468z M407.738,407.738H104.262V104.262   h303.475V407.738z"></path></g></svg>`;

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
        listSection: $('#list-section'),
        todoSection: $('#todo-section'),
    
        // temporary containers
        shadowTodo: document.createElement('div'),
        shadowDone: document.createElement('div'),

        init(ws) {
            this.ws = ws;

            // Reload page on navigation
            window.addEventListener('popstate', () => location.reload())

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
    
            // Update the title on focus loss
            this.title.addEventListener('focusout', evt => {
                this.onChangeTitle(this.title.innerText)
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

        // Hide the element passed
        hide(element) {
            element.classList.add("hidden");
        },

        // Show the element passed
        show(element) {
            element.classList.remove("hidden");
        },

        remove(element, immediate) {
            var exists = $(element)
            if (exists) {
                exists.classList.add("todo_removed");
                immediate ? exists.remove() : setTimeout(_ => exists.remove(), 1000);
            }
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
            localStorage.setItem(UI.listId ? `theme-${UI.listId}` : 'theme', color)
        },
    
        // Called on each received item
        onItemReceived(item) {
            switch (item.method) {
                case 'Ping':
                    console.log('connection openned')
                    this.enable()
                    break;
    
                case 'DeleteItem':
                    UI.remove(`#item-${item.itemId}`)
                    break;
    
                case 'AddItem':
                    UI.remove(`#item-${item.itemId}`, true)

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

        onListReceived(list) {
            switch (list.method) {
                case 'DeleteList':
                    UI.remove(`#list-${list.listId.substr(1)}`)
                    break;
    
                case "AddList":
                    UI.remove(`#list-${list.listId.substr(1)}`, true)
                    let node = document.createElement('div')
                    node.innerHTML = this.listTemplate(list)
                    node.querySelector('.todo__delete').addEventListener('click', e => WS.onDeleteList(e, list.listId))
                    const listNode = document.querySelector("#todolist-list")
                    listNode.appendChild(node)
                    break;
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
                ${IconDelete}
            </label>`;
        },

        listTemplate({ listId }) {
            const name = localStorage.getItem(listId)
            const linkText = !name || name === "" ? listId : name
            return `<li id="list-${listId.substr(1)}"><a href="/${listId}">${linkText}</a>${IconDelete}</li>`
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
            link.href = canvas.toDataURL('image/x-icon');
    
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
                data.forEach(val => {
                    if (val.method === "AddList" || val.method === "DeleteList") {
                        this.ui.onListReceived(val)
                        return
                    }
                    this.ui.onItemReceived(val)
                    this.ui.commitItemReceived();
                })
            } catch (e) {
                console.error(e)
            }
        },
    
        onOpen() {
            console.log('connection opened')
            if (!location.hash) {
                WS.fetchLists()
            } else {
                this.sock.send(JSON.stringify({
                    method: 'GetItems',
                    listId: this.ui.listId,
                }))
            }
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
        onDeleteList(event, listId) {
            setTimeout(() => {
                debugger;
                this.sock.send(JSON.stringify({
                    method: 'DeleteList',
                    listId,
                }))
            }, 400)
        },
        fetchLists() {
            setTimeout(() => {
                this.sock.send(JSON.stringify({
                    method: 'GetLists',
                }))
            }, 400)
        }
    }
    
    // UI & transport
    if (!location.hash) {
        // no hash, show a list of existing lists as well as link to create new
        const newListElem = document.querySelector("#new-list")
        newListElem.href = `/#${uuid()}:Untitled`;
        UI.hide(UI.todoSection)
        UI.show(UI.listSection)

        UI.init(WS)
        WS.init(UI)
        WS.fetchLists()
    } else {
        // app bootstrap
        UI.hide(UI.listSection)
        UI.show(UI.todoSection)

        hashAndTitle = location.hash.split(':')
        UI.title.innerText = decodeURIComponent(hashAndTitle[1] ? hashAndTitle[1] : hashAndTitle[0])
        UI.listId = hashAndTitle[0]

        UI.init(WS)
        WS.init(UI)
    }

    // Theme
    var color = localStorage.getItem(UI.listId ? `theme-${UI.listId}` : 'theme')
    if (color) UI.onChangeTheme(color)

    // List title
    var customTitle = localStorage.getItem(UI.listId)
    if (customTitle) UI.onChangeTitle(customTitle);
    else title.click()
    
})(window)
    
