require('react');
require('react-dom');
require('react-dom/server');
require('react-router');
require('history');
require('react-redux');
require('redux');

// add other libraries here as well
require('moment');
require('reselect');
require('react-custom-scrollbars');
require('react-lazyload');
require('react-select');
require('redux-undo');
require('whatwg-fetch');
require('babel-polyfill');
require('tsmonad');
require('lodash/mergeWith');
require('lodash/unionWith');

// re-bulma
require('insert-css');
require('re-bulma/lib/components/menu/menu');
require('re-bulma/lib/components/menu/menu-label');
require('re-bulma/lib/components/menu/menu-link');
require('re-bulma/lib/components/menu/menu-list');
require('re-bulma/lib/components/nav/nav');
require('re-bulma/lib/components/nav/nav-group');
require('re-bulma/lib/components/nav/nav-item');
require('re-bulma/lib/components/pagination/pagination');
require('re-bulma/lib/components/pagination/page-button');
require('re-bulma/lib/components/panel/panel');
require('re-bulma/lib/components/panel/panel-block');
require('re-bulma/lib/components/panel/panel-heading');
require('re-bulma/lib/elements/tag');
require('re-bulma/lib/elements/title');
require('re-bulma/lib/forms/input');
require('re-bulma/lib/grid/column');
require('re-bulma/lib/grid/columns');
require('re-bulma/lib/layout/container');
require('re-bulma/lib/layout/footer');
require('re-bulma/lib/layout/hero');
require('re-bulma/lib/layout/hero-body');
require('re-bulma/lib/layout/hero-head');
require('re-bulma/lib/layout/section');


// for dev
if (process.env.NODE_ENV !== 'production') {
    require('react-hot-loader');
    require('react-hot-loader/patch');
    require('webpack-dev-server/client');
    // require('webpack/hot/only-dev-server');
    require('redux-devtools');
    require('redux-devtools-log-monitor');
    require('redux-devtools-dock-monitor');
}