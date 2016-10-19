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

// require('classnames');
require('flexboxgrid');
require('react-flexbox-grid');

require('react-tap-event-plugin');
require('material-ui/styles/MuiThemeProvider');
require('material-ui/styles/getMuiTheme');
require('material-ui/AppBar');
require('material-ui/Drawer');
require('material-ui/IconButton');
require('material-ui/Divider');
require('material-ui/Paper');
require('material-ui/Table');
require('material-ui/FlatButton');
require('material-ui/RaisedButton');
require('material-ui/Dialog');
require('material-ui/Card');
require('material-ui/MenuItem');
require('material-ui/TextField');
require('material-ui/DatePicker');
require('material-ui/svg-icons/file/file-download');

// for dev
if (process.env.NODE_ENV !== 'production') {
    require('webpack-dev-server/client');
    // require('webpack/hot/only-dev-server');
    require('react-hot-loader');
    require('redux-devtools');
    require('redux-devtools-log-monitor');
    require('redux-devtools-dock-monitor');
}