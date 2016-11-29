import * as React from 'react';
import { Router, browserHistory } from 'react-router';
import { Provider } from 'react-redux';

import createRoutes from './routes';
import { configureStore, setAsCurrentStore } from '../store/configureStore';

export default class App extends React.Component<void, void> {
    render() {
        const store = configureStore(window['--app-initial']);
        setAsCurrentStore(store);

        return (
            <Provider store={store}>
                <Router history={browserHistory}>
                    {createRoutes({ store, first: { time: true } })}
                </Router>
            </Provider>
        );
    }
}

