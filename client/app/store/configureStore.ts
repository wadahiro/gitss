// Use DefinePlugin (Webpack) or loose-envify (Browserify)
// together with Uglify to strip the dev branch in prod build.
let cs;
if (process.env.NODE_ENV === 'production') {
    cs = require('./configureStore.prod').default;
} else {
    cs = require('./configureStore.dev').default;
}
export let configureStore = cs;

export let store = null;
export function getStore() { return store; }
export function setAsCurrentStore(s) {
    store = s;
    if (process.env.NODE_ENV !== 'production'
        && typeof window !== 'undefined') {
        window['store'] = store;
    }
}