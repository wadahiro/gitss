import { Store } from 'redux';
import { RootState } from '../reducers';

// Use DefinePlugin (Webpack) or loose-envify (Browserify)
// together with Uglify to strip the dev branch in prod build.
let cs;
if (process.env.NODE_ENV === 'production') {
    cs = require('./configureStore.prod').default;
} else {
    cs = require('./configureStore.dev').default;
}
export let configureStore: (initialState: any) => Store<RootState> = cs;

export let store: Store<RootState> = null;
export function getStore(): Store<RootState> { return store; }
export function setAsCurrentStore(s: Store<RootState>) {
    store = s;
    if (process.env.NODE_ENV !== 'production'
        && typeof window !== 'undefined') {
        window['store'] = store;
    }
}