import * as React from 'react';
import AppBar from 'material-ui/AppBar';

import { Toolbar, ToolbarGroup, ToolbarSeparator, ToolbarTitle } from 'material-ui/Toolbar';
import TextField from 'material-ui/TextField';
import Paper from 'material-ui/Paper';
import SearchIcon from 'material-ui/svg-icons/action/search';
import { List, ListItem } from 'material-ui/List';

import { Grid, Row, Col } from '../components/Grid';

export function NavBar(props) {
    const navStyle = {
        position: "fixed",
        backgroundColor: 'rgb(63, 81, 181)',
        width: '100%',
        zIndex: 1100,
        paddingLeft: 24,
        paddingRight: 24,
        top: 0,
        marginBottom: 80,
        boxSizing: 'border-box'
    };

    const style = {
        height: 50,
        width: 500,
        marginTop: 10,
        marginBottom: 10,
        marginLeft: 'auto',
        marginRight: 'auto',
        padding: 0
    };
    const listStyle = {
        height: '100%',
        width: '100%',
        paddingTop: 0,
        // padding: 0,
    };
    const listItemStyle = {
        height: '100%',
        width: '80%',
        paddingTop: 0,
        // padding: 0,
    };
    return (
        <div style={navStyle}>
            <Paper style={style} zDepth={4}>
                <List style={listStyle}>
                    <ListItem style={listItemStyle} primaryText={
                        <TextField
                            type='search'
                            hintText='Search'
                            underlineShow={false}
                            fullWidth={true}
                            onKeyDown={props.onKeyDown} />
                    } disabled leftIcon={<SearchIcon />}>
                    </ListItem>
                </List>
            </Paper>
        </div>
    );
}
