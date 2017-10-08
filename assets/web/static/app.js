const episodepublishform = {
    template: '<div><h3>Publish Episode</h3><form enctype="multipart/form-data" action="/admin/publish" method="post"><label for="title">Episode Title</label><input type="text" id="title" name="title"><label for="description">Episode Description</label><textarea name="description" id="description" cols="100" rows="20" style="resize: none;"></textarea><label for="file">Media File</label><input type="file" id="file" name="file"><label for="date">Publish Date</label><input type="date" id="date" name="date"><input type="submit" value="Publish"></form></div>'
}

const episodemanagement = {
    template: '<div><table style="width:100%"><tr><th>Title</th><th>URL</th><th>Actions</th></tr><tr v-for="item in items"><td>{{ item.title }}</td><td>{{ item.url }}</td><td><router-link :to="\'edit/\' + item.id">Edit</router-link></td></tr></table></div>',
    data() {
        return {
            loading: false,
            items: null,
            error: null
        }
    },
    created() {
        // fetch the data when the view is created and the data is
        // already being observed
        this.fetchData()
    },
    watch: {
        // call again the method if the route changes
        '$route': 'fetchData'
    },
    methods: {
        fetchData() {
            this.error = this.items = []
            this.loading = true

            getEpisodes((err, items) => {
                this.loading = false
                if (err) {
                    this.error = err.toString()
                } else {
                	console.log(items);
                    var t = JSON.parse(items).items
                    for (var i = t.length - 1; i >= 0; i--) {
                        this.items.push({
                            title: t[i].title,
                            url: t[i].url,
                            id: t[i].id
                        })
                    }
                }
            })
        }
    }
}

const episodeedit = {
	template: '<div><div><h3>Edit Episode</h3><form enctype="multipart/form-data" action="/admin/edit" method="post"><label for="title">Episode Title</label><input type="text" id="title" name="title"><label for="description">Episode Description</label><textarea name="description" id="description" cols="100" rows="20" style="resize: none;"></textarea><label for="date">Publish Date</label><input type="date" id="date" name="date"><input type="submit" value="Publish"></form></div></div>',
	data() {
        return {
            loading: false,
            episode: null,
            error: null
        }
    },
    created() {
        // fetch the data when the view is created and the data is
        // already being observed
        this.fetchData()
    },
    watch: {
        // call again the method if the route changes
        '$route': 'fetchData'
    },
    methods: {
        fetchData() {
            this.error = this.items = []
            this.loading = true

            getEpisodes((err, items) => {
                this.loading = false
                if (err) {
                    this.error = err.toString()
                } else {
                    var t = JSON.parse(items).items
                    for (var i = t.length - 1; i >= 0; i--) {
                    	if (t[i].id == route.params.id) 
                        this.items.push({
                            title: t[i].title,
                            url: t[i].url
                        })
                    }
                }
                console.log(this.items)
            })
        }
    }
}

const customcss = {
    template: '<div><h3>Edit CSS</h3><form action="/admin/css" method="post" enctype="multipart/form-data"><label for="css">Custom CSS</label><textarea name="css" id="css" cols="120" rows="20">{{ css }}</textarea><br /><input type="submit" value="Submit"></form></div>',
    data() {
        return {
            loading: false,
            css: null,
            error: null
        }
    },
    created() {
        // fetch the data when the view is created and the data is
        // already being observed
        this.fetchData()
    },
    watch: {
        // call again the method if the route changes
        '$route': 'fetchData'
    },
    methods: {
        fetchData() {
            this.error = this.css = null
            this.loading = true

            getCss((err, css) => {
                this.loading = false
                if (err) {
                    this.error = err.toString()
                } else {
                    this.css = css
                }
            })
        }
    }
}

const routes = [
    { path: '/publish', component: episodepublishform },
    { path: '/manage', component: episodemanagement },
    { path: '/theme', component: customcss },
    { path: '/edit/:id', component: episodeedit}
]

const router = new VueRouter({
    routes
})

const app = new Vue({
    router,
    data: {
        header: 'Pogo Admin',
    }
}).$mount('#app')

function getCss(callback) {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.onreadystatechange = function() {
        if (xmlHttp.readyState == 4 && xmlHttp.status == 200)
            callback(null, xmlHttp.responseText)
    }
    xmlHttp.open("GET", "/admin/css", true);
    xmlHttp.send(null);
}

function getEpisodes(callback) {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.onreadystatechange = function() {
        if (xmlHttp.readyState == 4 && xmlHttp.status == 200)
            callback(null, xmlHttp.responseText)
    }
    xmlHttp.open("GET", "/json", true);
    xmlHttp.send(null);
}