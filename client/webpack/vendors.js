require('react');
require('react-dom');
require('react-dom/server');
require('react-router');
require('react-redux');
require('redux');

// add other libraries here as well
require('moment');
require('reselect');
require('redux-undo');
require('whatwg-fetch');
require('babel-polyfill');

// for dev
if (process.env.NODE_ENV !== 'production') {
    require('webpack-dev-server/client');
    // require('webpack/hot/only-dev-server');
    require('react-hot-loader');
    require('redux-devtools');
    require('redux-devtools-log-monitor');
    require('redux-devtools-dock-monitor');
}