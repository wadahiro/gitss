import * as React from 'react';
import * as ReactDOM from 'react-dom';
import { AppContainer } from 'react-hot-loader';
import { Provider } from 'react-redux'
import * as ReactDOMServer from 'react-dom/server';
import { match, RouterContext } from 'react-router'

require('font-awesome/css/font-awesome.css');
const insertCss = require('insert-css');
const css = require('re-bulma/build/css');

insertCss(css, { prepend: true });

// init promise polyfill
window['Promise'] = window['Promise'] || Promise;
// init fetch polyfill
// window.self = window;
require('whatwg-fetch');

const App = require('./app/router/App').default;

ReactDOM.render(
    <AppContainer>
        <App />
    </AppContainer>,
    document.getElementById('app')
);

if (module['hot']) {
    module['hot'].accept('./app/router/App', () => {
        const NextApp = require('./app/router/App').default;
        ReactDOM.render(
            <AppContainer>
                <NextApp />
            </AppContainer>,
            document.getElementById('app')
        );
    });
}