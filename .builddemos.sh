# Assumes that .runtests.sh has been ran prior to execution.

blacklist="engo.io/engo/demos/demoutils,engo.io/engo/demos/tilemap"

for dir in `pwd`/demos/*/
do
    # Formatting the directory to be usable by Go
    # In the form: engo.io/engo/demos/demoutils
    dir=${dir%*/}
    dir=${dir#$GOPATH/src/}

    # Ignore the directory if it's in the blacklist
    if [[ $blacklist == *"${dir}"* ]]
    then
        echo "Skipping ${dir}"
        continue
    fi

    local_demo_dir=${dir#engo.io/engo/}
    local_assets_dir=$local_demo_dir'/assets'
    site_dir='site/'$dir

    go_file=$local_demo_dir'/'$(basename $local_demo_dir)'.go'
    js_file=$dir'/script.js'

    echo 'site/'$js_file
    html_file=$site_dir'/index.html'

    mkdir -p $site_dir
    gopherjs build -o 'site/'$js_file $dir

    assets=$GOPATH/src/$dir/assets
    if [ -d "$assets" ]; then
        cp -r $assets $site_dir
    fi
    
    data=$GOPATH/src/$dir/data
    if [ -d "$data" ]; then
        cp -r $data $site_dir
    fi

    script_tag="<script src='/$js_file'></script"
    code_tag='<pre>'$(cat $go_file)'</pre>'

    echo $script_tag $code_tag > $html_file
done
