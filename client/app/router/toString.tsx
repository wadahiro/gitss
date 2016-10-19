import * as React from 'react';
import { Provider } from 'react-redux';
import { renderToString } from 'react-dom/server';
import { match, RouterContext } from 'react-router';
import * as Helmet from 'react-helmet';

import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import { indigo500, indigo700, indigo400 } from 'material-ui/styles/colors';

import createRoutes from './routes';
import { configureStore, setAsCurrentStore } from '../store/configureStore';


/**
 * Handle HTTP request at Golang server
 *
 * @param   {Object}   options  request options
 * @param   {Function} cbk      response callback
 */
export default function (options, cbk) {

    let result = {
        uuid: options.uuid,
        app: null,
        title: null,
        meta: null,
        initial: null,
        error: null,
        redirect: null
    };

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

    const store = configureStore();
    setAsCurrentStore(store);

    try {
        match({ routes: createRoutes({ store, first: { time: false } }), location: options.url }, (error, redirectLocation, renderProps) => {
            try {
                if (error) {
                    result.error = error;

                } else if (redirectLocation) {
                    result.redirect = redirectLocation.pathname + redirectLocation.search;

                } else {
                    result.app = renderToString(
                        <div>hoge</div>
                    );
                    const { title, meta } = Helmet.rewind();
                    result.title = title.toString();
                    result.meta = meta.toString();
                    result.initial = JSON.stringify(store.getState());
                }
            } catch (e) {
                result.error = e;
            }
            return cbk(JSON.stringify(result));
        });
    } catch (e) {
        result.error = e;
        return cbk(JSON.stringify(result));
    }
}