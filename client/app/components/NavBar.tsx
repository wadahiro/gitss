import * as React from 'react';
import AppBar from 'material-ui/AppBar';

import { Toolbar, ToolbarGroup, ToolbarSeparator, ToolbarTitle } from 'material-ui/Toolbar';
import TextField from 'material-ui/TextField';
import Paper from 'material-ui/Paper';
import SearchIcon from 'material-ui/svg-icons/action/search';
import { List, ListItem } from 'material-ui/List';

import { Grid, Row, Col } from '../components/Grid';

export function NavBar(props) {
    const styles = {
        position: "fixed",
        top: 0
    };

    const style = {
        height: 50,
        width: 500,
        margin: 20,
        padding: 0,
        textAlign: 'center',
        display: 'inline-block',
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
        <div style={{ marginBottom: 40 }}>
            <AppBar
                style={styles}
                iconElementLeft={
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
                } />

            <Toolbar>
                <ToolbarGroup firstChild={true}>
                    <Paper style={style} zDepth={1} >
                        <TextField
                            type='search'
                            floatingLabelText='Search string...'
                            fullWidth={true}
                            />
                    </Paper>
                </ToolbarGroup>
                <ToolbarGroup>
                    <ToolbarTitle text="Options" />

                </ToolbarGroup>
            </Toolbar>
        </div >
    );
}
