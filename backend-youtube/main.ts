import { Innertube, UniversalCache } from "youtubei.js";
import { MusicInlineBadge } from "youtubei.js/dist/src/parser/nodes";
import { Hono } from "hono";

const innertube = await Innertube.create({
    lang: 'en',
    location: 'US',
    visitor_data: '',
    po_token: '',
    retrieve_player: true,
    enable_safety_mode: false,
    generate_session_locally: false,
    enable_session_cache: true,
    device_category: 'desktop',
    cookie: '',
    cache: new UniversalCache(
        true,
        './cache'
    )
});

const app = new Hono();

app.get('/search', async (c) => {
   const query = c.req.query('query');
   if (query === '' || query === undefined) {
       c.status(400);
       return c.json({
          error: 'empty query',
       });
   }

   console.log(`Received search request for query: ${query}`);

    const search = await innertube.music.search(query, {
        type: 'song'
    });

    const result = search.songs.contents.slice(0, 5).map(song => {
        const artists = song.artists?.map(x => x.name).join(', ');
        const durationSec = song.duration?.seconds!;
        const explicit = song.badges?.find(item => {
            const badge = item as MusicInlineBadge;
            return badge.icon_type === 'MUSIC_EXPLICIT_BADGE';
        }) !== undefined;

        return {
            id: song.id,
            title: song.title,
            authors: artists,
            thumbnail: song.thumbnail!.contents[song.thumbnails!.length - 1].url,
            length: durationSec,
            explicit: explicit
        }
    });

    return c.json(result);
});

app.get('/download', async (c) => {
   const id = c.req.query('id');
   if (id === '' || id === undefined) {
       c.status(400);
       return c.json({
           error: 'empty id',
       });
   }

    console.log(`Received download request for video ID: ${id}`);

    const info = await innertube.music.getInfo(id);
    const durationSec = info.basic_info.duration!;
    const stream = await info.download({
        type: 'audio',
        quality: 'best',
        format: 'mp4'
    });
    const chunks: Uint8Array[] = [];
    for await (const chunk of stream) {
        chunks.push(chunk);
    }
    const buffer = Buffer.concat(chunks);
    const rawFile = `./dl/${id}.mp4`;
    await Bun.write(rawFile, buffer);
    const fixedFile = `./dl/${id}.m4a`;
    const ffmpeg = Bun.spawn({
        cmd: [
            'ffmpeg',
            '-y',
            '-i', rawFile,
            '-acodec', 'copy',
            '-vn',
            '-t', durationSec.toString(),
            fixedFile,
        ],
        stdout: 'inherit',
        stderr: 'inherit',
    });

    const code = await ffmpeg.exited;
    if (code === 0) {
        console.log(`→ Fixed audio saved to ${fixedFile}`);
    } else {
        console.error('ffmpeg process failed');
    }

    return c.json({
        status: 'success',
    });
});

export default app