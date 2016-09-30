import { Action } from 'redux';

const { ActionCreators } = require('redux-undo');

export type Actions =
    SearchItem |
    AddItem |
    DeleteItem |
    ModifyQuantity |
    RestoreSavedHistory |
    ModifyMetadata
    ;

export interface SearchItem extends Action {
    type: 'SEARCH_ITEM';
    payload: {
        searchWord: string;
    }
}
export function searchItem(searchWord: string): SearchItem {
    return {
        type: 'SEARCH_ITEM',
        payload: {
            searchWord
        }
    };
}

export interface AddItem extends Action {
    type: 'ADD_ITEM';
    payload: {
        itemId: string;
    }
}
export function addItem(itemId: string): AddItem {
    return {
        type: 'ADD_ITEM',
        payload: {
            itemId
        }
    };
}

export interface DeleteItem extends Action {
    type: 'DELETE_ITEM';
    payload: {
        itemId: string;
    }
}
export function deleteItem(itemId: string): DeleteItem {
    return {
        type: 'DELETE_ITEM',
        payload: {
            itemId
        }
    };
}

export interface ModifyQuantity extends Action {
    type: 'MOD_QUANTITY';
    payload: {
        itemId: string;
        quantity: number;
    };
}
export function modifyQuantity(itemId: string, quantity: number): ModifyQuantity {
    return {
        type: 'MOD_QUANTITY',
        payload: {
            itemId,
            quantity
        }
    };
}

export interface RestoreSavedHistory extends Action {
    type: 'RESTORE_SAVED_HISTORY';
    payload: {
        date: string;
    }
}
export function restoreSavedHistory(date: string): RestoreSavedHistory {
    return {
        type: 'RESTORE_SAVED_HISTORY',
        payload: {
            date
        }
    };
}

export interface ModifyMetadata extends Action {
    type: 'MOD_METADATA';
    payload: {
        value: {
            [index: string]: string
        };
    };
}
export function modifyMetadata(value: { [index: string]: string }): ModifyMetadata {
    return {
        type: 'MOD_METADATA',
        payload: {
            value
        }
    };
}

export function undo() {
    return ActionCreators.undo();
}

export function redo() {
    return ActionCreators.redo();
}