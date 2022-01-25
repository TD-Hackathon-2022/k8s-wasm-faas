use std::collections::HashMap;
use std::fs::OpenOptions;
use std::io::Write;

use crate::lambda::lambda_handler;

mod lambda;

fn main() {
    let mut params_map = HashMap::new();
    let mut result_out_put_file_path = String::new();
    for arg in std::env::args() {
        result_out_put_file_path = arg.clone();
        if arg.contains("=") {
            let index = arg.find("=");
            match index {
                None => {}
                Some(_) => {
                    let mut name = String::new();
                    let mut value = String::new();
                    let mut value_flag = false;
                    for c in arg.chars() {
                        if !value_flag {
                            if c == '=' {
                                value_flag = true;
                                continue;
                            }
                            name = name + &String::from(c)
                        } else {
                            value = value + &String::from(c)
                        }
                    }
                    params_map.insert(name, value);
                }
            }
        }
    }

    let result = lambda_handler(params_map);

    let mut file = OpenOptions::new()
        .read(true).write(true).open(result_out_put_file_path).unwrap();
    file.write(result.to_string().as_bytes()).unwrap();
}
