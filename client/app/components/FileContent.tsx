import * as React from 'react';

import Chip from 'material-ui/Chip';
import { indigo50 } from 'material-ui/styles/colors';

import { FileMetadata, Highlight } from '../reducers';

const CRLF = /\r\n|\r|\n/g;

interface Props {
    metadata: FileMetadata[];
    contents: Highlight[];
}

interface State {
    highlight: any;
}

export class FileContent extends React.Component<Props, State>{

    highlight(offset: number, highlightNums: number[], content: string, ext: string) {
        if (window['EnlighterJS']) {
            const pre = document.createElement('pre');
            pre.setAttribute('data-enlighter-lineoffset', String(offset));
            pre.setAttribute('data-enlighter-highlight', highlightNums.join(','));
            pre.appendChild(document.createTextNode(content));

            const div = document.createElement('div');

            const enlighter = new window['EnlighterJS'](pre, {
                indent: 2
            }, div);
            enlighter.enlight(true);

            return <div dangerouslySetInnerHTML={{ __html: div.outerHTML.replace('/\n/g', '') }} />;
        }

        return null;
    }

    render() {
        const { metadata, contents} = this.props;

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
                {contents.map((content, i) => {

                    // stript
                    const codeList = content.content.split(CRLF);

                    // calc highlight lines
                    const highlightNums = codeList.reduce((s, x, index) => {
                        if (x.search(/@GITK_MARK_(PRE|POST)@/g) !== -1) {
                            s.push(index + content.offset);
                        }
                        return s;
                    }, []);

                    // format
                    const code = codeList.join('\n')
                        .replace(/@GITK_MARK_(PRE|POST)@/g, '');

                    return (
                        <div key={content.offset}>
                            {this.highlight(content.offset, highlightNums, code, metadata[0].ext)}
                        </div>
                    );
                })}
            </div>
        );
    }
}