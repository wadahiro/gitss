import * as React from 'react';
import { render } from 'react-dom';
import { Router, browserHistory } from 'react-router';
import { Provider } from 'react-redux';

import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import { indigo500, indigo700, indigo400 } from 'material-ui/styles/colors';

import toString from './toString';
// import { Promise } from 'when';
import createRoutes from './routes';
import { configureStore, setAsCurrentStore } from '../store/configureStore';

export function run() {
    // init promise polyfill
    window['Promise'] = window['Promise'] || Promise;
    // init fetch polyfill
    // window.self = window;
    require('whatwg-fetch');

    // Needed for onTouchTap
    // http://stackoverflow.com/a/34015469/988941
    const injectTapEventPlugin = require('react-tap-event-plugin');
    injectTapEventPlugin();

    const muiTheme = getMuiTheme({
        fontFamily: 'Helvetica,Arial,sans-serif',
        tableRow: {
            height: 30
        },
        tableHeaderColumn: {
            height: 30
        },
        palette: {
            primary1Color: indigo500,
            primary2Color: indigo700
        }
    });

    const store = configureStore(window['--app-initial']);
    setAsCurrentStore(store);

    render(
        <Provider store={store}>
            <MuiThemeProvider muiTheme={muiTheme}>
                <Router history={browserHistory}>{createRoutes({ store, first: { time: true } })}</Router>
            </MuiThemeProvider>
        </Provider>,
        document.getElementById('app')
    );

}

// Export it to render on the Golang sever, keep the name sync with -
// https://github.com/olebedev/go-starter-kit/blob/master/src/app/server/react.go#L65
export const renderToString = toString;

// require('../css');

// Style live reloading
// if (module['hot']) {
//     let c = 0;
//     module['hot'].accept('../css', () => {
//         require('../css');
//         const a = document.createElement('a');
//         const link = document.querySelector('link[rel="stylesheet"]');
//         a.href = link['href'];
//         a.search = '?' + c++;
//         link['href'] = a.href;
//     });
// }