library IEEE;

use IEEE.std_logic_1164.all;

entity xorg2 is

  port(
    i_A : in  std_logic;
    i_B : in  std_logic;
    o_F : out std_logic
    );

end xorg2;

architecture dataflow of xorg2 is
begin

  o_F <= i_A xor i_B;

end dataflow;
