import * as React from 'react';

import Chip from 'material-ui/Chip';
import { indigo50 } from 'material-ui/styles/colors';

import { FileMetadata } from '../reducers';

interface Props {
    metadata: FileMetadata[];
    content: string;
}

export class FileContent extends React.Component<Props, void>{
    componentDidMount() {
        const pre = this.refs['code']
        if (pre && window['EnlighterJS']) {
            pre['enlight']({
                language: 'js',
                indent: 2
            });
        }
    }
    render() {
        const { metadata, content} = this.props;
        const styles = {
            chip: {
                margin: 4,
            },
            wrapper: {
                display: 'flex',
                flexWrap: 'wrap',
            },
        };
        return (
            <div>
                <div style={styles.wrapper}>
                    {metadata.map(x => {
                        return (
                            <Chip
                                key={`${x.project}_${x.repo}_${x.refs}_${x.path}`}
                                backgroundColor={indigo50}
                                style={styles.chip}
                                >
                                {`${x.project}/${x.repo} (${x.refs}) - ${x.path}`}
                            </Chip>
                        );
                    })}
                </div>
                <pre ref='code'>
                    {content}
                </pre >
            </div>
        );
    }
}