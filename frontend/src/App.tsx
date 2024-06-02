import { useState } from 'react'
import './App.css'
import axios from 'axios'

function App() {
  const [url, setURL] = useState<string>("")
  const [shortURL, setShortURL] = useState<string>("")

  const generateURL = async () => {
    if (urlValidator(url) != null) {
      const res = await axios.post("http://localhost:8000/api/generate", { url, })
      setShortURL("http://localhost:8000/" + res.data.short_url)
    } else {

      alert("Invalid URL.")
    }
  }

  const urlValidator = (url: string): RegExpMatchArray | null => {
    return url.match(/(http(s)?:\/\/.)?(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)/g)
  }

  return (
    <main className='w-screen h-screen flex flex-col gap-6 items-center pt-16'>
      <h1 className='text-3xl font-bold'><span className=' text-green-500'>Shortener</span> Link</h1>
      <div className="flex flex-col p-6 shadow-sm border border-green-500 rounded-md gap-4 w-[40%]">
        <div className='flex justify-center'>
          <input
            type="url"
            className=' bg-white p-3 w-[30rem] text-black'
            placeholder='www.example.com'
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setURL(e.target.value)}
          />
          <button className='bg-green-500 px-4 font-bold hover:bg-green-700' onClick={generateURL}>Generate</button>
        </div>
        <div>
          <p className='text-center'><span className=' text-green-500 font-bold'>Shortener</span> Link is a free tool to shorten URLs and generate short links <br />
            URL shortener allows to create a shortened link making it easy to share</p>
        </div>
      </div>
      <div className='flex flex-col p-6 shadow-sm border border-green-500 rounded-md items-center gap-4 w-[40%] justify-center'>
        <div className='flex items-center gap-3'>
          <h1 className=' text-green-500 font-bold text-lg'>OUTPUT</h1>
          <span>{shortURL.length == 0 ? "-" : shortURL}</span>
        </div>
        <div className='flex gap-3'>
          <button className='bg-green-500 p-3 font-bold hover:bg-green-700 rounded-md'>Copy</button>
          <button
            className='bg-blue-500 p-3 font-bold hover:bg-blue-700 rounded-md'
            disabled={shortURL.length == 0}
            onClick={(e: React.MouseEvent<HTMLButtonElement>) => {
              window.location.href = shortURL
            }}
          >Go to link</button>
        </div>
      </div>
    </main>
  )
}

export default App
